// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package memo

import (
	memo3 "github.com/heyjun3/dforget/backend/app/memo"
	memo2 "github.com/heyjun3/dforget/backend/domain/memo"
	"github.com/heyjun3/dforget/backend/infra/memo"
	"github.com/uptrace/bun"
)

// Injectors from wire.go:

func InitializeMemoHandler(db *bun.DB) *MemoHandler {
	memoRepository := memo.NewMemoRepository(db)
	registerMemoService := memo2.NewRegisterMemoService(memoRepository)
	memoUsecase := memo3.NewMemoUsecase(registerMemoService, memoRepository)
	memoHandler := NewMemoHandler(memoRepository, registerMemoService, memoUsecase)
	return memoHandler
}
