package main

import (
	"log/slog"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	cfg "github.com/heyjun3/dforget/backend/config"
	"github.com/heyjun3/dforget/backend/presentation"
)

func main() {
	mux := presentation.NewServer(cfg.NewConfig())

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "https://*.trycloudflare.com"},
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
