//go:build wireinject
// +build wireinject

package server

import (
	"database/sql"
	"net/http"

	"github.com/google/wire"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func initializeMemoHandler(db *bun.DB) *MemoHandler {
	wire.Build(
		NewMemoRepository,
		NewMemoHandler,
	)
	return nil
}

func initializeOIDCHandler(conf Config) *OIDCHandler {
	wire.Build(
		provideHttpClient,
		NewOIDCHandler,
	)
	return nil
}

func provideHttpClient() *http.Client {
	return &http.Client{}
}

type DBConfigIF interface {
	DBDSN() string
}

func provideOpenDB(conf DBConfigIF) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.DBDSN())))
	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}
func InitDBConn(conf DBConfigIF) *bun.DB {
	wire.Build(
		provideOpenDB,
	)
	return nil
}
