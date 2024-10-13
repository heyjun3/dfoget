package memo

import (
	"github.com/google/uuid"
)

type Memo struct {
	ID     uuid.UUID
	UserId uuid.UUID
	Title  string
	Text   string
}

func NewMemo(title, text string, userId uuid.UUID, opts ...Option) (*Memo, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return nil, err
	}
	memo := &Memo{
		ID:     id,
		UserId: userId,
		Title:  title,
		Text:   text,
	}
	for _, opt := range opts {
		memo, err = opt(memo)
		if err != nil {
			return nil, err
		}
	}
	return memo, nil
}

type Option func(*Memo) (*Memo, error)

func WithID(id string) Option {
	return func(memo *Memo) (*Memo, error) {
		ID, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		memo.ID = ID
		return memo, nil
	}
}
func WithUUID(id uuid.UUID) Option {
	return func(memo *Memo) (*Memo, error) {
		memo.ID = id
		return memo, nil
	}
}

func (m Memo) IsEqual(memo Memo) bool {
	return m.Title == memo.Title && m.Text == memo.Text
}
