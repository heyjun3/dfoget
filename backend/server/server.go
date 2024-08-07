package server

import (
	"net/http"

	"github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
)

func New(conf Config) *http.ServeMux {
	mux := http.NewServeMux()
	db := InitDBConn(conf)
	memo := initializeMemoHandler(db)
	path, handler := memov1connect.NewMemoServiceHandler(memo)
	mux.Handle(path, handler)
	return mux
}
