package server

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
)

func New(conf Config) *http.ServeMux {
	mux := http.NewServeMux()
	db := InitDBConn(conf)
	memo := initializeMemoHandler(db)
	interceptors := connect.WithInterceptors(NewAuthInterceptor(conf))
	path, handler := memov1connect.NewMemoServiceHandler(memo, interceptors)
	mux.Handle(path, handler)

	oidcHandler := initializeOIDCHandler(conf)
	mux.HandleFunc("GET /oidc", oidcHandler.recieveRedirect)

	return mux
}
