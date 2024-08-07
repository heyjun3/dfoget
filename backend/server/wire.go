//go:build wireinject
// +build wireinject

package server

import (
	"database/sql"

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
