package server_test

import (
	"context"
	"errors"
	"fmt"
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
	conf := server.NewConfig(server.WithDBName("test"))
	ResetModel(server.InitDBConn(conf))
	mux := server.New(conf)
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

	t.Run("invalid id", func(t *testing.T) {
		id := "id"
		res, err := client.RegisterMemo(context.Background(),
			connect.NewRequest(&memov1.RegisterMemoRequest{
				Memo: &memov1.Memo{
					Id:    server.Ptr(id),
					Title: "test",
					Text:  "test",
				},
			}),
		)

		assert.Nil(t, res)
		assert.Error(t, err)
		fmt.Println(connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			fmt.Println(connectErr.Message())
			fmt.Println(connectErr.Details())
		}
	})
}
