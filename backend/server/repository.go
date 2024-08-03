package server

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Memo struct {
	ID    uuid.UUID
	Title string
	Text  string
}

type MemoRepository struct {
	db *bun.DB
}

func NewMemoRepository(db *bun.DB) *MemoRepository {
	return &MemoRepository{
		db: db,
	}
}

func (r *MemoRepository) Save(ctx context.Context, memos []Memo) ([]Memo, error) {
	dm := memoToDM(memos)
	_, err := r.db.NewInsert().Model(&dm).
		On("CONFLICT (id) DO UPDATE").
		Set(strings.Join([]string{
			"title = EXCLUDED.title",
			"text = EXCLUDED.text",
		}, ",")).
		Exec(ctx)
	return memos, err
}

func (r *MemoRepository) Find(ctx context.Context) ([]Memo, error) {
	dm := make([]MemoDM, 0)
	err := r.db.NewSelect().Model(&dm).Scan(ctx)
	return dmToMemo(dm), err
}

func (r *MemoRepository) DeleteByIds(ctx context.Context, ids []uuid.UUID) ([]uuid.UUID, error) {
	_, err := r.db.NewDelete().Model((*MemoDM)(nil)).
		Where("id IN (?)", bun.In(ids)).Exec(ctx)
	return ids, err
}

func dmToMemo(memoDM []MemoDM) []Memo {
	memos := make([]Memo, 0, len(memoDM))
	for _, dm := range memoDM {
		memos = append(memos, Memo{
			ID:    dm.ID,
			Title: dm.Title,
			Text:  dm.Text,
		})
	}
	return memos
}

func memoToDM(memos []Memo) []MemoDM {
	dm := make([]MemoDM, 0, len(memos))
	for _, memo := range memos {
		dm = append(dm, MemoDM{
			ID:    memo.ID,
			Title: memo.Title,
			Text:  memo.Text,
		})
	}
	return dm
}
