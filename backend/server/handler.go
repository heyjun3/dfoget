package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"connectrpc.com/connect"
	memoapp "github.com/heyjun3/dforget/backend/app/memo"
	"github.com/heyjun3/dforget/backend/domain/memo"
	memov1 "github.com/heyjun3/dforget/backend/gen/api/memo/v1"
	memov1connect "github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
	"github.com/heyjun3/dforget/backend/lib"
)

func Ptr[T any](v T) *T {
	return &v
}

var _ memov1connect.MemoServiceHandler = (*MemoHandler)(nil)

func NewMemov1Memo(memo *memo.Memo) *memov1.Memo {
	return &memov1.Memo{
		Id:    Ptr(memo.ID.String()),
		Title: memo.Title,
		Text:  memo.Text,
	}
}
func NewMemov1Memos(memos []*memo.Memo) []*memov1.Memo {
	dto := make([]*memov1.Memo, 0, len(memos))
	for _, memo := range memos {
		dto = append(dto, NewMemov1Memo(memo))
	}
	return dto
}

func NewMemoHandler(memoRepository *MemoRepository,
	registerMemoService *memo.RegisterMemoService,
	memoUsecase *memoapp.MemoUsecase,
) *MemoHandler {
	return &MemoHandler{
		memoRepository:      memoRepository,
		registerMemoService: registerMemoService,
		memoUsecase:         memoUsecase,
	}
}

type MemoHandler struct {
	memoRepository      *MemoRepository
	registerMemoService *memo.RegisterMemoService
	memoUsecase         *memoapp.MemoUsecase
}

func (h MemoHandler) RegisterMemo(ctx context.Context, req *connect.Request[memov1.RegisterMemoRequest]) (
	*connect.Response[memov1.RegisterMemoResponse], error,
) {
	id := req.Msg.Memo.Id
	title := req.Msg.Memo.Title
	text := req.Msg.Memo.Text
	memo, err := h.memoUsecase.RegisterMemo(ctx, id, title, text)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(
		&memov1.RegisterMemoResponse{
			Memo: NewMemov1Memo(memo),
		},
	)
	return res, nil
}

func (h MemoHandler) GetMemo(ctx context.Context, req *connect.Request[memov1.GetMemoRequest]) (
	*connect.Response[memov1.GetMemoResponse], error,
) {
	memos, err := h.memoUsecase.FindMemos(ctx)
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(
		&memov1.GetMemoResponse{
			Memo: NewMemov1Memos(memos),
		},
	)
	return res, nil
}

func (h MemoHandler) GetMemoServerStream(
	ctx context.Context,
	req *connect.Request[memov1.GetMemoServerStreamRequest],
	stream *connect.ServerStream[memov1.GetMemoServerStreamResponse],
) error {
	var prevMemo []*memo.Memo
	for {
		memos, err := h.memoUsecase.FindMemos(ctx)
		if err != nil {
			return fmt.Errorf("failed get memos")
		}
		if len(prevMemo) != 0 {
			if prevMemo[0].IsEqual(*memos[0]) {
				time.Sleep(time.Second * 10)
				continue
			}
		}
		if err := stream.Send(&memov1.GetMemoServerStreamResponse{
			Memo: NewMemov1Memos(memos),
		}); err != nil {
			return err
		}
		prevMemo = memos
		time.Sleep(time.Second * 10)
	}
}

func (h MemoHandler) DeleteMemo(ctx context.Context, req *connect.Request[memov1.DeleteMemoRequest]) (
	*connect.Response[memov1.DeleteMemoResponse], error,
) {
	err := h.memoUsecase.DeleteMemo(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(
		&memov1.DeleteMemoResponse{
			Id: req.Msg.Id,
		},
	)
	return res, nil
}

func (h MemoHandler) MemoStream(ctx context.Context,
	stream *connect.BidiStream[memov1.MemoStreamRequest, memov1.MemoStreamResponse]) error {
	for {
		msg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to receive request: %w", err))
		}

		sub, err := lib.GetSubValue(ctx)
		if err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}

		id := msg.Memo.Id
		title := msg.Memo.Title
		text := msg.Memo.Text
		memo, err := h.registerMemoService.Execute(
			ctx, sub, id, title, text)
		if err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
		if err := stream.Send(&memov1.MemoStreamResponse{
			Memo: &memov1.Memo{
				Id:    Ptr(memo.ID.String()),
				Title: memo.Title,
				Text:  memo.Text,
			},
		}); err != nil {
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to send response: %w", err))
		}
	}
}

type OIDCHandler struct {
	conf       Config
	httpClient httpClient
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewOIDCHandler(conf Config, client httpClient) *OIDCHandler {
	return &OIDCHandler{
		conf:       conf,
		httpClient: client,
	}
}

func (h OIDCHandler) recieveRedirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "recieve oidc redirect")
	code := r.URL.Query().Get("code")
	if code == "" {
		slog.ErrorContext(ctx, "code is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	formData := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {h.conf.oidc.redirectUri},
		"client_id":     {h.conf.oidc.clientId},
		"client_secret": {h.conf.oidc.clientSecret},
	}
	req, err := http.NewRequest("POST", h.conf.oidc.tokenUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := h.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var token OIDCToken
	if err := json.Unmarshal(buf, &token); err != nil {
		slog.ErrorContext(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:     AuthCookieName,
		Value:    token.IdToken,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	slog.InfoContext(ctx, "oidc verified")
	http.Redirect(w, r, h.conf.frontEndURL, http.StatusTemporaryRedirect)
}

type OIDCToken struct {
	IdToken string `json:"id_token"`
}
