// go:build wireinject
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

	memoapp "github.com/heyjun3/dforget/backend/app/memo"
	cfg "github.com/heyjun3/dforget/backend/config"
	"github.com/heyjun3/dforget/backend/domain/memo"
)

func InitializeMemoHandler(db *bun.DB) *MemoHandler {
	wire.Build(
		NewMemoRepository,
		memo.NewRegisterMemoService,
		memoapp.NewMemoUsecase,
		NewMemoHandler,
		wire.Bind(new(memo.MemoRepositoryInterface), new(*MemoRepository)),
		wire.Bind(new(memoapp.MemoRepositoryInterface), new(*MemoRepository)),
	)
	return nil
}

func InitializeOIDCHandler(conf cfg.Config) *OIDCHandler {
	wire.Build(
		provideHttpClient,
		NewOIDCHandler,
	)
	return nil
}

func provideHttpClient() httpClient {
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
