package main

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/heyjun3/dforget/backend/server"
	"github.com/rs/cors"
)

func main() {
	mux := server.New(server.NewConfig())

	http.ListenAndServe(
		"dev:8080",
		cors.AllowAll().Handler((h2c.NewHandler(mux, &http2.Server{}))),
	)
}
