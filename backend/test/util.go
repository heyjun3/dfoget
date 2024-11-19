package test

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/heyjun3/dforget/backend/server"
)

func ResetModel(db *bun.DB) {
	models := []interface{}{
		server.MemoDM{},
	}
	for _, model := range models {
		db.NewDelete().Model(&model).Exec(context.Background())
	}
}