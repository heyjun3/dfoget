package chat

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/heyjun3/dforget/backend/lib"
	"github.com/stretchr/testify/assert"
)

func TestChat(t *testing.T) {
	t.Run("test add message", func(t *testing.T) {
		ctx := lib.SetSubKey(context.Background(), uuid.New().String())
		room, err := newRoom("test room")
		assert.NoError(t, err)

		room.AddMessage(ctx, "test message")

		assert.GreaterOrEqual(t, 1, len(room.Messages))
		assert.Equal(t, "test message", room.Messages[0].text)
	})

	t.Run("delete message", func(t *testing.T) {
		ctx := lib.SetSubKey(context.Background(), uuid.New().String())
		room, err := newRoom("test room")
		assert.NoError(t, err)

		room.AddMessage(ctx, "test message 1")
		room.AddMessage(ctx, "test message 2")

		deleteId := room.Messages[0].id
		room.DeleteMessage(ctx, deleteId)

		assert.Equal(t, 1, len(room.Messages))
		assert.Equal(t, "test message 2", room.Messages[0].text)
	})
}
