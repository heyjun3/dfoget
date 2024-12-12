// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"database/sql"
	memo2 "github.com/heyjun3/dforget/backend/app/memo"
	"github.com/heyjun3/dforget/backend/config"
	"github.com/heyjun3/dforget/backend/domain/memo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"net/http"
)

// Injectors from wire.go:

func InitializeMemoHandler(db *bun.DB) *MemoHandler {
	memoRepository := NewMemoRepository(db)
	registerMemoService := memo.NewRegisterMemoService(memoRepository)
	memoUsecase := memo2.NewMemoUsecase(registerMemoService, memoRepository)
	memoHandler := NewMemoHandler(memoRepository, registerMemoService, memoUsecase)
	return memoHandler
}

func InitializeOIDCHandler(conf server.Config) *OIDCHandler {
	serverHttpClient := provideHttpClient()
	oidcHandler := NewOIDCHandler(conf, serverHttpClient)
	return oidcHandler
}

func InitDBConn(conf DBConfigIF) *bun.DB {
	db := provideOpenDB(conf)
	return db
}

// wire.go:

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
