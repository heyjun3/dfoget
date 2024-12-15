package memo_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"github.com/heyjun3/dforget/backend/domain/memo"
	memodm "github.com/heyjun3/dforget/backend/infra/memo"
	"github.com/heyjun3/dforget/backend/lib"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func TestMemoRepository(t *testing.T) {
	dsn := "postgres://dev:dev@postgres:5432/test?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	repo := memodm.NewMemoRepository(db)
	db.NewTruncateTable().
		Model((*memodm.MemoDM)(nil)).Exec(context.Background())

	t.Run("save, find, and delete memo", func(t *testing.T) {
		id, _ := uuid.NewV7()
		userId := uuid.New()
		ctx := lib.SetSubKey(context.Background(), userId.String())
		memos := []*memo.Memo{
			{ID: id, Title: "title", Text: "text", UserId: userId},
		}

		_, err := repo.Save(ctx, memos)
		assert.NoError(t, err)

		memos, err = repo.Find(ctx)
		assert.NoError(t, err)
		assert.Equal(t, []*memo.Memo{
			{
				ID:     id,
				Title:  "title",
				Text:   "text",
				UserId: userId,
			}}, memos)

		m, err := repo.GetById(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, memos[0], m)

		err = repo.DeleteByIds(ctx, []uuid.UUID{id})
		assert.NoError(t, err)

		memos, err = repo.Find(ctx)
		assert.NoError(t, err)
		assert.Equal(t, []*memo.Memo{}, memos)
	})
}
