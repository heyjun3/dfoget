package presentation

import (
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"

	cfg "github.com/heyjun3/dforget/backend/config"
	"github.com/heyjun3/dforget/backend/gen/api/chat/v1/chatv1connect"
	"github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
	"github.com/heyjun3/dforget/backend/lib/database"
	"github.com/heyjun3/dforget/backend/presentation/auth"
	"github.com/heyjun3/dforget/backend/presentation/chat"
	"github.com/heyjun3/dforget/backend/presentation/memo"
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

func NewServer(conf cfg.Config) *http.ServeMux {
	mux := http.NewServeMux()
	db := database.InitDBConn(conf)
	interceptors := connect.WithInterceptors(auth.NewAuthInterceptorV2(conf))

	memo := memo.InitializeMemoHandler(db)
	path, handler := memov1connect.NewMemoServiceHandler(memo, interceptors)
	mux.Handle(path, loggerMiddleware(handler))

	chat := chat.InitChatHandler(db)
	path, handler = chatv1connect.NewChatServiceHandler(chat, interceptors)
	mux.Handle(path, loggerMiddleware(handler))

	oidcHandler := auth.InitializeOIDCHandler(conf)
	mux.HandleFunc("GET /oidc", oidcHandler.RecieveRedirect)

	return mux
}
