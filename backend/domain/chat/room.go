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
	id        uuid.UUID
	name      string
	messages  []Message
	createdAt time.Time
}

type RoomWithoutMessage struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
}

type RoomOption func(r *Room) *Room

func WithReconstructRoom(id uuid.UUID, messages []Message, createdAt time.Time) RoomOption {
	return func(r *Room) *Room {
		r.id = id
		r.messages = messages
		r.createdAt = createdAt
		return r
	}
}

func NewRoom(name string, opts ...RoomOption) (*Room, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("room name is required")
	}
	r := &Room{
		id:        id,
		name:      name,
		createdAt: time.Now(),
	}
	for _, opt := range opts {
		r = opt(r)
	}
	return r, nil
}

func (r *Room) AddMessage(ctx context.Context, text string) error {
	userId, err := lib.GetSubValue(ctx)
	if err != nil {
		return err
	}
	message, err := NewMessage(userId, r.id, text)
	if err != nil {
		return err
	}
	r.messages = append(r.messages, *message)
	return nil
}

func (r *Room) DeleteMessage(ctx context.Context, messageId uuid.UUID) error {
	userId, err := lib.GetSubValue(ctx)
	if err != nil {
		return err
	}
	r.messages = slices.DeleteFunc(r.messages, func(elem Message) bool {
		return elem.id.String() == messageId.String() && elem.userID == userId
	})
	return nil
}

func (r *Room) Get() (id uuid.UUID, name string, messages []Message, createdAt time.Time) {
	return r.id, r.name, r.messages, r.createdAt
}

type Message struct {
	id        uuid.UUID
	userID    uuid.UUID
	roomID    uuid.UUID
	text      string
	createdAt time.Time
}

type MessageOption func(*Message) *Message

func WithReconstruct(id uuid.UUID, createdAt time.Time) MessageOption {
	return func(m *Message) *Message {
		m.id = id
		m.createdAt = createdAt
		return m
	}
}

func NewMessage(userID, roomID uuid.UUID, text string, opts ...MessageOption) (*Message, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	m := &Message{
		id:        id,
		userID:    userID,
		roomID:    roomID,
		text:      text,
		createdAt: time.Now(),
	}
	for _, opt := range opts {
		m = opt(m)
	}
	return m, nil
}

func (m *Message) Get() (id, userID, roomID uuid.UUID, text string, createdAt time.Time) {
	return m.id, m.userID, m.roomID, m.text, m.createdAt
}
