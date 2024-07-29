package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	memov1 "github.com/heyjun3/dforget/backend/gen/api/memo/v1"
	"github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
	"github.com/heyjun3/dforget/backend/server"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func TestGetMemo(t *testing.T) {
	mux := server.New()
	srv := httptest.NewServer(h2c.NewHandler(mux, &http2.Server{}))
	defer srv.Close()

	client := memov1connect.NewMemoServiceClient(
		http.DefaultClient,
		srv.URL,
	)

	t.Run("run get memo", func(t *testing.T) {
		expect := []*memov1.Memo{
			{Id: server.Ptr("test"), Title: "test title", Text: "text"},
		}

		res, err := client.GetMemo(context.Background(),
			connect.NewRequest(&memov1.GetMemoRequest{}))

		assert.NoError(t, err)
		assert.Equal(t, expect, res.Msg.Memo)
	})
}
