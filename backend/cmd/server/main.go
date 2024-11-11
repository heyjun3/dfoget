package main

import (
	"log/slog"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/heyjun3/dforget/backend/server"
	"github.com/rs/cors"
)

func main() {
	mux := server.New(server.NewConfig())

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{
			http.MethodOptions,
			http.MethodPost,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	slog.Info("start listen. port: 8080")
	http.ListenAndServe(
		"dev:8080",
		c.Handler((h2c.NewHandler(mux, &http2.Server{}))),
	)
}
