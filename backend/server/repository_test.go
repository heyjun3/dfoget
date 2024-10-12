package server_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/heyjun3/dforget/backend/server"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func TestMemoRepository(t *testing.T) {
	dsn := "postgres://dev:dev@postgres:5432/test?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	repo := server.NewMemoRepository(db)
	db.NewTruncateTable().
		Model((*server.MemoDM)(nil)).Exec(context.Background())

	t.Run("save, find, and delete memo", func(t *testing.T) {
		id, _ := uuid.NewV7()
		userId := uuid.New()
		ctx := server.SetSubKey(context.Background(), userId.String())
		memos := []*server.Memo{
			{ID: id, Title: "title", Text: "text", UserId: userId},
		}

		_, err := repo.Save(ctx, memos)
		assert.NoError(t, err)

		memos, err = repo.Find(ctx, userId)
		assert.NoError(t, err)
		assert.Equal(t, []*server.Memo{
			{
				ID:     id,
				Title:  "title",
				Text:   "text",
				UserId: userId,
			}}, memos)

		memo, err := repo.GetById(ctx, userId, id)
		assert.NoError(t, err)
		assert.Equal(t, memos[0], memo)

		_, err = repo.DeleteByIds(ctx,
			userId, []uuid.UUID{id})
		assert.NoError(t, err)

		memos, err = repo.Find(context.Background(), userId)
		assert.NoError(t, err)
		assert.Equal(t, []*server.Memo{}, memos)
	})
}
