package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	memov1 "github.com/heyjun3/dforget/backend/gen/api/memo/v1"
	"github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
	"github.com/heyjun3/dforget/backend/server"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func TestGetMemo(t *testing.T) {
	ResetModel(server.OpenDB(server.NewConfig().TESTDBDSN()))
	mux := server.New(server.NewConfig().TESTDBDSN())
	srv := httptest.NewServer(h2c.NewHandler(mux, &http2.Server{}))
	defer srv.Close()

	client := memov1connect.NewMemoServiceClient(
		http.DefaultClient,
		srv.URL,
	)

	t.Run("run register memo", func(t *testing.T) {
		id, err := uuid.NewV7()
		assert.NoError(t, err)

		res, err := client.RegisterMemo(context.Background(),
			connect.NewRequest(&memov1.RegisterMemoRequest{
				Memo: &memov1.Memo{
					Id:    server.Ptr(id.String()),
					Title: "test",
					Text:  "test",
				},
			}),
		)

		expect := &memov1.Memo{
			Id:    server.Ptr(id.String()),
			Title: "test",
			Text:  "test",
		}
		assert.NoError(t, err)
		assert.Equal(t, expect, res.Msg.Memo)

		getres, err := client.GetMemo(context.Background(),
			connect.NewRequest(&memov1.GetMemoRequest{}))
		assert.NoError(t, err)
		assert.Equal(t, []*memov1.Memo{expect}, getres.Msg.Memo)

		deleteres, err := client.DeleteMemo(context.Background(),
			connect.NewRequest(&memov1.DeleteMemoRequest{
				Id: []string{id.String()},
			}),
		)
		assert.NoError(t, err)
		assert.Equal(t, []string{id.String()}, deleteres.Msg.Id)

		getres, err = client.GetMemo(context.Background(),
			connect.NewRequest(&memov1.GetMemoRequest{}))
		assert.NoError(t, err)
		assert.Equal(t, []*memov1.Memo(nil), getres.Msg.Memo)
	})
}
