//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
	"github.com/uptrace/bun"
)

func initializeMemoHandler(db *bun.DB) *MemoHandler {
	wire.Build(
		NewMemoRepository,
		NewMemoHandler,
	)
	return nil
}
