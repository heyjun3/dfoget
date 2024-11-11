package server

import (
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
)

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New()
		slog.InfoContext(ctx, "start request", "requestID", id)
		next.ServeHTTP(w, r)
		slog.InfoContext(ctx, "done request", "requestID", id)
	})
}

func New(conf Config) *http.ServeMux {
	mux := http.NewServeMux()
	db := InitDBConn(conf)
	memo := initializeMemoHandler(db)
	interceptors := connect.WithInterceptors(NewAuthInterceptorV2(conf))
	path, handler := memov1connect.NewMemoServiceHandler(memo, interceptors)
	mux.Handle(path, loggerMiddleware(handler))

	oidcHandler := initializeOIDCHandler(conf)
	mux.HandleFunc("GET /oidc", oidcHandler.recieveRedirect)

	return mux
}
