package chat

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"database/sql"

	"github.com/google/uuid"
	"github.com/heyjun3/dforget/backend/lib"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunslog"
)

func SetupTestDBConnection() *bun.DB {
	dsn := "postgres://dev:dev@postgres:5432/test?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bunslog.NewQueryHook(
		bunslog.WithQueryLogLevel(slog.LevelDebug),
	))
	return db
}

func ResetDB(db *bun.DB, models []interface{}) {
	for _, model := range models {
		_, err := db.NewDelete().Model(model).Where("1 = 1").Exec(context.Background())
		if err != nil {
			panic(err)
		}
	}
}

func TestChatRepository(t *testing.T) {
	db := SetupTestDBConnection()
	ResetDB(db, []interface{}{(*MessageDM)(nil), (*RoomDM)(nil)})

	repo := NewChatRepository(db)
	roomID, _ := uuid.NewV7()
	messageID1, _ := uuid.NewV7()
	messageID2, _ := uuid.NewV7()
	userId, _ := uuid.NewV7()
	ctx := lib.SetSubKey(context.Background(), userId.String())

	messages := []*MessageDM{
		{
			ID:        messageID1,
			UserID:    userId,
			RoomID:    roomID,
			Text:      "message 1",
			CreatedAt: time.Now().In(time.UTC).Round(time.Millisecond),
		},
		{
			ID:        messageID2,
			UserID:    userId,
			RoomID:    roomID,
			Text:      "message 2",
			CreatedAt: time.Now().In(time.UTC).Round(time.Millisecond),
		},
	}
	room := RoomDM{
		ID:        roomID,
		Name:      "test room",
		CreatedAt: time.Now().In(time.UTC).Round(time.Microsecond),
		Messages:  (messages),
	}
	t.Run("save room", func(t *testing.T) {
		err := repo.Save(ctx, dmToRoom(&room))

		assert.NoError(t, err)
	})
	t.Run("exists", func(t *testing.T) {
		exist, err := repo.Exists(ctx, "test room")

		assert.NoError(t, err)
		assert.True(t, exist)

		exist, err = repo.Exists(ctx, "test room 1")

		assert.NoError(t, err)
		assert.False(t, exist)
	})

	t.Run("get room", func(t *testing.T) {
		r, err := repo.GetRoom(ctx, roomID)

		assert.NoError(t, err)
		assert.NotNil(t, r)
		dm := roomToDM(r)
		assert.Equal(t, &room, dm)
	})

	t.Run("get room without messages", func(t *testing.T) {
		rooms, err := repo.GetRoomsWithoutMessage(ctx)

		assert.NoError(t, err)
		assert.Equal(t, room.ID, rooms[0].ID)
		assert.Equal(t, room.Name, rooms[0].Name)
		assert.Equal(t, room.CreatedAt, rooms[0].CreatedAt)
	})

	t.Run("delete", func(t *testing.T) {
		err := repo.DeleteById(ctx, roomID)

		assert.NoError(t, err)

		exist, err := repo.Exists(ctx, "test room")

		assert.NoError(t, err)
		assert.False(t, exist)
	})
}
