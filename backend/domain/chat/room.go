package chat

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/heyjun3/dforget/backend/lib"
)

type Room struct {
	ID        uuid.UUID
	Name      string
	Messages  []Message
	CreatedAt time.Time
}

type RoomWithoutMessage struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
}

func newRoom(name string) (*Room, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("room name is required")
	}
	return &Room{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
}

func (r *Room) AddMessage(ctx context.Context, text string) error {
	userId, err := lib.GetSubValue(ctx)
	if err != nil {
		return err
	}
	message, err := newMessage(userId, r.ID, text)
	if err != nil {
		return err
	}
	r.Messages = append(r.Messages, *message)
	return nil
}

func (r *Room) DeleteMessage(ctx context.Context, messageId uuid.UUID) error {
	userId, err := lib.GetSubValue(ctx)
	if err != nil {
		return err
	}
	r.Messages = slices.DeleteFunc(r.Messages, func(elem Message) bool {
		return elem.id.String() == messageId.String() && elem.userID == userId
	})
	return nil
}

type Message struct {
	id        uuid.UUID
	userID    uuid.UUID
	roomID    uuid.UUID
	text      string
	createdAt time.Time
}

func newMessage(userID, roomID uuid.UUID, text string) (*Message, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &Message{
		id:        id,
		userID:    userID,
		roomID:    roomID,
		text:      text,
		createdAt: time.Now(),
	}, nil
}
