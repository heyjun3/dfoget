package memo_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	cfg "github.com/heyjun3/dforget/backend/config"
	memov1 "github.com/heyjun3/dforget/backend/gen/api/memo/v1"
	"github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
	"github.com/heyjun3/dforget/backend/presentation"
	"github.com/heyjun3/dforget/backend/server"
	"github.com/heyjun3/dforget/backend/test"
)

func TestMemoHandler(t *testing.T) {
	conf := cfg.NewConfig(
		cfg.WithDBName("test"),
		cfg.WithPubKey(test.PublicKey),
	)
	test.ResetModel(server.InitDBConn(conf))
	mux := presentation.NewServer(conf)
	srv := httptest.NewServer(h2c.NewHandler(mux, &http2.Server{}))
	defer srv.Close()

	interceptors := connect.WithInterceptors(test.NewSetCookieInterceptor())
	client := memov1connect.NewMemoServiceClient(
		http.DefaultClient,
		srv.URL,
		interceptors,
	)

	t.Run("run register memo", func(t *testing.T) {
		res, err := client.RegisterMemo(context.Background(),
			connect.NewRequest(&memov1.RegisterMemoRequest{
				Memo: &memov1.Memo{
					Title: "test title",
					Text:  "test text",
				},
			}),
		)

		assert.NoError(t, err)
		assert.NotNil(t, res.Msg.Memo.Id)
		assert.Equal(t, "test title", res.Msg.Memo.Title)
		assert.Equal(t, "test text", res.Msg.Memo.Text)

		getres, err := client.GetMemo(context.Background(),
			connect.NewRequest(&memov1.GetMemoRequest{}))
		assert.NoError(t, err)
		assert.NotNil(t, getres.Msg.Memo[0].Id)
		assert.Equal(t, "test title", getres.Msg.Memo[0].Title)
		assert.Equal(t, "test text", getres.Msg.Memo[0].Text)

		res, err = client.RegisterMemo(context.Background(),
			connect.NewRequest(&memov1.RegisterMemoRequest{
				Memo: &memov1.Memo{
					Id:    getres.Msg.Memo[0].Id,
					Title: "test title v2",
					Text:  "test text v2",
				},
			}),
		)
		assert.NoError(t, err)
		assert.Equal(t, "test title v2", res.Msg.Memo.Title)
		assert.Equal(t, "test text v2", res.Msg.Memo.Text)

		getres, err = client.GetMemo(context.Background(),
			connect.NewRequest(&memov1.GetMemoRequest{}))
		assert.NoError(t, err)
		assert.Equal(t, "test title v2", getres.Msg.Memo[0].Title)
		assert.Equal(t, "test text v2", getres.Msg.Memo[0].Text)

		deleteres, err := client.DeleteMemo(context.Background(),
			connect.NewRequest(&memov1.DeleteMemoRequest{
				Id: []string{*res.Msg.Memo.Id},
			}),
		)
		assert.NoError(t, err)
		assert.Equal(t, []string{*res.Msg.Memo.Id}, deleteres.Msg.Id)

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
