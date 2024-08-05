package server_test

import (
	"context"

	"github.com/heyjun3/dforget/backend/server"
	"github.com/uptrace/bun"
)

func ResetModel(db *bun.DB) {
	models := []interface{}{
		server.MemoDM{},
	}
	for _, model := range models {
		db.NewTruncateTable().Model(&model).Exec(context.Background())
	}
}
