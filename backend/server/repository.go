package server

import (
	"context"
	"strings"

	"github.com/google/uuid"
	memoapp "github.com/heyjun3/dforget/backend/app/memo"
	"github.com/heyjun3/dforget/backend/domain/memo"
	"github.com/heyjun3/dforget/backend/lib"
	"github.com/uptrace/bun"
)

var _ memo.MemoRepositoryInterface = (*MemoRepository)(nil)
var _ memoapp.MemoRepositoryInterface = (*MemoRepository)(nil)

type MemoRepository struct {
	db *bun.DB
}

func NewMemoRepository(db *bun.DB) *MemoRepository {
	return &MemoRepository{
		db: db,
	}
}

func (r *MemoRepository) Save(ctx context.Context, memos []*memo.Memo) (
	[]*memo.Memo, error) {
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
	ctx context.Context) ([]*memo.Memo, error) {
	userId, err := lib.GetSubValue(ctx)
	if err != nil {
		return nil, err
	}
	dm := make([]MemoDM, 0)
	err = r.db.NewSelect().
		Model(&dm).
		Where("user_id = ?", userId.String()).
		Scan(ctx)
	return dmToMemos(dm), err
}

func (r *MemoRepository) GetById(
	ctx context.Context, id uuid.UUID) (
	*memo.Memo, error) {
	var dm MemoDM
	userId, err := lib.GetSubValue(ctx)
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

func dmToMemo(memoDM MemoDM) *memo.Memo {
	return &memo.Memo{
		ID:     memoDM.ID,
		UserId: memoDM.UserId,
		Title:  memoDM.Title,
		Text:   memoDM.Text,
	}
}

func dmToMemos(memoDM []MemoDM) []*memo.Memo {
	memos := make([]*memo.Memo, 0, len(memoDM))
	for _, dm := range memoDM {
		memos = append(memos, dmToMemo(dm))
	}
	return memos
}

func memoToDM(memos []*memo.Memo) []MemoDM {
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
