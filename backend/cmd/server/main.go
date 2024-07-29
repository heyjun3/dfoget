package main

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
	"github.com/heyjun3/dforget/backend/server"
)

func main() {
	memo := &server.MemoHandler{}
	mux := http.NewServeMux()

	path, handler := memov1connect.NewMemoServiceHandler(memo)
	mux.Handle(path, handler)


	http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
