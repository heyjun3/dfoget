package test_test

import (
	"net/http/httptest"
	"testing"

	"github.com/heyjun3/dforget/backend/server"
	"github.com/heyjun3/dforget/backend/test"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func TestChatHandler(t *testing.T) {
	conf := server.NewConfig(server.WithDBName("test"), server.WithPubKey(publicKey))
	test.ResetModel(server.InitDBConn(conf))
	mux := server.New(conf)
	srv := httptest.NewServer(h2c.NewHandler(mux, &http2.Server{}))
	defer srv.Close()

}