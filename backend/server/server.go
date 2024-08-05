package server

import (
	"database/sql"
	"net/http"

	"github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func OpenDB(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}

func New(dsn string) *http.ServeMux {
	mux := http.NewServeMux()
	db := OpenDB(dsn)
	memo := initializeMemoHandler(db)
	path, handler := memov1connect.NewMemoServiceHandler(memo)
	mux.Handle(path, handler)
	return mux
}
