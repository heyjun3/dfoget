package server

import (
	"context"

	"connectrpc.com/connect"
	memov1 "github.com/heyjun3/dforget/backend/gen/api/memo/v1"
)

func Ptr[T any](v T) *T {
	return &v
}

type MemoHandler struct{}

func (h MemoHandler) RegisterMemo(ctx context.Context, req *connect.Request[memov1.RegisterMemoRequest]) (
	*connect.Response[memov1.RegisterMemoResponse], error,
) {
	id := "id"
	res := connect.NewResponse(
		&memov1.RegisterMemoResponse{
			Memo: &memov1.Memo{
				Id:    &id,
				Title: "test",
				Text:  "test",
			},
		},
	)
	return res, nil
}

func (h MemoHandler) GetMemo(ctx context.Context, req *connect.Request[memov1.GetMemoRequest]) (
	*connect.Response[memov1.GetMemoResponse], error,
) {
	res := connect.NewResponse(
		&memov1.GetMemoResponse{
			Memo: []*memov1.Memo{
				{Id: Ptr("test"), Title: "test title", Text: "text"},
			},
		},
	)
	return res, nil
}

func (h MemoHandler) DeleteMemo(ctx context.Context, req *connect.Request[memov1.DeleteMemoRequest]) (
	*connect.Response[memov1.DeleteMemoResponse], error,
) {
	res := connect.NewResponse(
		&memov1.DeleteMemoResponse{
			Id: []string{},
		},
	)
	return res, nil
}
