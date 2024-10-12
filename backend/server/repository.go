package server

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ MemoRepositoryInterface = (*MemoRepository)(nil)

type MemoRepository struct {
	db *bun.DB
}

func NewMemoRepository(db *bun.DB) *MemoRepository {
	return &MemoRepository{
		db: db,
	}
}

func (r *MemoRepository) Save(ctx context.Context, memos []*Memo) (
	[]*Memo, error) {
	dm := memoToDM(memos)
	_, err := r.db.NewInsert().Model(&dm).
		On("CONFLICT (id) DO UPDATE").
		Set(strings.Join([]string{
			"title = EXCLUDED.title",
			"text = EXCLUDED.text",
			"user_id = EXCLUDED.user_id",
		}, ",")).
		Exec(ctx)
	return memos, err
}

func (r *MemoRepository) Find(
	ctx context.Context, userId uuid.UUID) ([]*Memo, error) {
	dm := make([]MemoDM, 0)
	err := r.db.NewSelect().
		Model(&dm).
		Where("user_id = ?", userId.String()).
		Scan(ctx)
	return dmToMemos(dm), err
}

func (r *MemoRepository) GetById(
	ctx context.Context, userId uuid.UUID, id uuid.UUID) (
	*Memo, error) {
	var dm MemoDM
	userId, err := GetSubValue(ctx)
	if err != nil {
		return nil, err
	}
	err = r.db.NewSelect().
		Model(&dm).
		Where("user_id = ?", userId.String()).
		Where("id = ?", id.String()).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return dmToMemo(dm), err
}

func (r *MemoRepository) DeleteByIds(
	ctx context.Context, userId uuid.UUID, ids []uuid.UUID) (
	[]uuid.UUID, error) {
	_, err := r.db.NewDelete().
		Model((*MemoDM)(nil)).
		Where("id IN (?)", bun.In(ids)).
		Where("user_id = ?", userId.String()).
		Exec(ctx)
	return ids, err
}

func dmToMemo(memoDM MemoDM) *Memo {
	return &Memo{
		ID:     memoDM.ID,
		UserId: memoDM.UserId,
		Title:  memoDM.Title,
		Text:   memoDM.Text,
	}
}

func dmToMemos(memoDM []MemoDM) []*Memo {
	memos := make([]*Memo, 0, len(memoDM))
	for _, dm := range memoDM {
		memos = append(memos, dmToMemo(dm))
	}
	return memos
}

func memoToDM(memos []*Memo) []MemoDM {
	dm := make([]MemoDM, 0, len(memos))
	for _, memo := range memos {
		dm = append(dm, MemoDM{
			ID:     memo.ID,
			UserId: memo.UserId,
			Title:  memo.Title,
			Text:   memo.Text,
		})
	}
	return dm
}
