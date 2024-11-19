package test_test

import (
	"testing"

	"github.com/heyjun3/dforget/backend/server"
	"github.com/heyjun3/dforget/backend/test"
)

func TestChatHandler(t *testing.T) {
	conf := server.NewConfig(server.WithDBName("test"), server.WithPubKey(publicKey))
	test.ResetModel(server.InitDBConn(conf))
	mux := server.New(conf)
	srv := httptest.NewServer(h2c.NewHandler(mux, &http2.Server{}))
	defer srv.Close()

}