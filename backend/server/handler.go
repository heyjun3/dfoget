package server

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	memov1 "github.com/heyjun3/dforget/backend/gen/api/memo/v1"
)

func Ptr[T any](v T) *T {
	return &v
}

func NewMemoHandler(memoRepository *MemoRepository) *MemoHandler {
	return &MemoHandler{
		memoRepository: memoRepository,
	}
}

type MemoHandler struct {
	memoRepository *MemoRepository
}

func (h MemoHandler) RegisterMemo(ctx context.Context, req *connect.Request[memov1.RegisterMemoRequest]) (
	*connect.Response[memov1.RegisterMemoResponse], error,
) {
	id := req.Msg.Memo.Id
	title := req.Msg.Memo.Title
	text := req.Msg.Memo.Text
	var opts []Option
	if id != nil {
		opts = append(opts, WithID(*id))
	}
	memo, err := NewMemo(title, text, opts...)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	_, err = h.memoRepository.Save(context.Background(), []Memo{*memo})
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(
		&memov1.RegisterMemoResponse{
			Memo: &memov1.Memo{
				Id:    Ptr(memo.ID.String()),
				Title: memo.Title,
				Text:  memo.Text,
			},
		},
	)
	return res, nil
}

func (h MemoHandler) GetMemo(ctx context.Context, req *connect.Request[memov1.GetMemoRequest]) (
	*connect.Response[memov1.GetMemoResponse], error,
) {
	memos, err := h.memoRepository.Find(context.Background())
	if err != nil {
		return nil, err
	}
	var memosDTO []*memov1.Memo
	for _, memo := range memos {
		memosDTO = append(memosDTO, &memov1.Memo{
			Id:    Ptr(memo.ID.String()),
			Title: memo.Title,
			Text:  memo.Text,
		})
	}
	res := connect.NewResponse(
		&memov1.GetMemoResponse{
			Memo: memosDTO,
		},
	)
	return res, nil
}

func (h MemoHandler) DeleteMemo(ctx context.Context, req *connect.Request[memov1.DeleteMemoRequest]) (
	*connect.Response[memov1.DeleteMemoResponse], error,
) {
	ids := req.Msg.Id
	var uuids []uuid.UUID
	for _, id := range ids {
		uu, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, uu)
	}
	_, err := h.memoRepository.DeleteByIds(context.Background(), uuids)
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

type OIDCHandler struct {
	conf       Config
	httpClient *http.Client
}

func NewOIDCHandler(conf Config, httpClient *http.Client) *OIDCHandler {
	return &OIDCHandler{
		conf:       conf,
		httpClient: httpClient,
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

	slog.InfoContext(ctx, "oidc verified")
	http.Redirect(w, r, h.conf.frontEndURL, http.StatusTemporaryRedirect)
}
