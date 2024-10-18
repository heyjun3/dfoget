package chat

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID
	Name      string
	Messages  []Message
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

func (r *Room) AddMessage(message Message) {
	r.Messages = append(r.Messages, message)
}

func (r *Room) DeleteMessage(messageId uuid.UUID) {
	r.Messages = slices.DeleteFunc(r.Messages, func(elem Message) bool {
		return elem.ID.String() == messageId.String()
	})
}

type Message struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	RoomID    uuid.UUID
	Text      string
	CreatedAt time.Time
}

func NewMessage(userID, roomID uuid.UUID, text string) (*Message, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &Message{
		ID:        id,
		UserID:    userID,
		RoomID:    roomID,
		Text:      text,
		CreatedAt: time.Now(),
	}, nil
}
