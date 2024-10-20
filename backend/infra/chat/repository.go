package chat

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/heyjun3/dforget/backend/domain/chat"
	"github.com/uptrace/bun"
)

type RoomDM struct {
	bun.BaseModel `bun:"table:rooms"`
	ID            uuid.UUID    `bun:"id,pk,type:uuid"`
	Name          string       `bun:"name,type:text,notnull,unique"`
	CreatedAt     time.Time    `bun:"type:timestamp,notnull,default:now()"`
	Messages      []*MessageDM `bun:"rel:has-many,join:id=room_id"`
}

type MessageDM struct {
	bun.BaseModel `bun:"table:messages"`
	ID            uuid.UUID `bun:"id,pk,type:uuid"`
	UserID        uuid.UUID `bun:"user_id,type:uuid,notnull"`
	RoomID        uuid.UUID `bun:"room_id,type:uuid,notnull"`
	Text          string    `bun:"type:text,notnull"`
	CreatedAt     time.Time `bun:"type:timestamp,notnull,default:now()"`
}

var _ chat.CreateRoomRepositoryInterface = (*ChatRepository)(nil)

type ChatRepository struct {
	db *bun.DB
}

func NewChatRepository(db *bun.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (r *ChatRepository) Exists(ctx context.Context, name string) (bool, error) {
	var room *RoomDM
	if err := r.db.NewSelect().Model(room).Where("name = ?", name).Scan(ctx); err != nil {
		return false, err
	}
	return room != nil, nil
}
