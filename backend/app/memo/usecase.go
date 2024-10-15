package memo

import (
	"context"

	"github.com/google/uuid"
	domain "github.com/heyjun3/dforget/backend/domain/memo"
	"github.com/heyjun3/dforget/backend/lib"
)

type MemoRepositoryInterface interface {
	Find(context.Context) ([]*domain.Memo, error)
	DeleteByIds(context.Context, []uuid.UUID) error
}

func NewMemoUsecase(
	registerMemoService *domain.RegisterMemoService,
	memoRepositoryInterface MemoRepositoryInterface,
) *MemoUsecase {
	return &MemoUsecase{
		registerMemoService:     registerMemoService,
		memoRepositoryInterface: memoRepositoryInterface,
	}
}

type MemoUsecase struct {
	registerMemoService     *domain.RegisterMemoService
	memoRepositoryInterface MemoRepositoryInterface
}

func (u *MemoUsecase) RegisterMemo(ctx context.Context, id *string, title, text string) (*domain.Memo, error) {
	sub, err := lib.GetSubValue(ctx)
	if err != nil {
		return nil, err
	}
	memo, err := u.registerMemoService.Execute(ctx, sub, id, title, text)
	if err != nil {
		return nil, err
	}
	return memo, nil
}

func (u *MemoUsecase) FindMemos(ctx context.Context) ([]*domain.Memo, error) {
	return u.memoRepositoryInterface.Find(ctx)
}

func (u *MemoUsecase) DeleteMemo(ctx context.Context, ids []string) error {
	uuids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		uu, err := uuid.Parse(id)
		if err != nil {
			return err
		}
		uuids = append(uuids, uu)
	}
	return u.memoRepositoryInterface.DeleteByIds(ctx, uuids)
}
