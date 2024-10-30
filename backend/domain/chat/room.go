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

type Message struct {
	id        uuid.UUID
	userID    uuid.UUID
	roomID    uuid.UUID
	text      string
	createdAt time.Time
}

type Identifiler interface {
	SetID(uuid.UUID)
}

func (m *Message) SetID(id uuid.UUID) {
	m.id = id
}

type Option[T Identifiler] func(v T) T
type MessageOption func(*Message) *Message

func WithID[T Identifiler](id uuid.UUID) Option[T] {
	return func(m T) T{
		m.SetID(id)
		return m
	}
}
func WithCreatedAt(createdAt time.Time) MessageOption {
	return func(m *Message) *Message {
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
