package server

import (
	"context"

	"connectrpc.com/connect"
	memov1 "github.com/heyjun3/dforget/backend/gen/api/memo/v1"
)

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
