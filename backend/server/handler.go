package server

import (
	"context"

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
		return nil, err
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
