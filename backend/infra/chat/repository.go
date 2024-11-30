package chat

import (
	"context"
	"time"

	"github.com/google/uuid"
	appchat "github.com/heyjun3/dforget/backend/app/chat"
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

var _ appchat.RoomRepository = (*ChatRepository)(nil)

type ChatRepository struct {
	db *bun.DB
}

func NewChatRepository(db *bun.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (r *ChatRepository) Save(ctx context.Context, room *chat.Room) error {
	dm := roomToDM(room)
	_, err := r.db.NewInsert().Model(dm).Exec(ctx)
	if err != nil {
		return err
	}
	if len(dm.Messages) > 0 {
		_, err = r.db.NewInsert().Model(&dm.Messages).Exec(ctx)
	}
	return err
}

func (r *ChatRepository) Exists(ctx context.Context, name string) (bool, error) {
	return r.db.NewSelect().Model((*RoomDM)(nil)).Where("name = ?", name).Exists(ctx)
}

func (r *ChatRepository) GetRoom(ctx context.Context, id uuid.UUID) (*chat.Room, error) {
	room := new(RoomDM)
	if err := r.db.NewSelect().
		Model(room).
		Relation("Messages").
		Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return dmToRoom(room), nil
}

func (r *ChatRepository) GetRoomsWithoutMessage(ctx context.Context) ([]*chat.RoomWithoutMessage, error) {
	var rooms []*RoomDM
	if err := r.db.NewSelect().Model(&rooms).Scan(ctx); err != nil {
		return nil, nil
	}
	return dmToRoomWithoutMessages(rooms), nil
}

func (r *ChatRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*RoomDM)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}
