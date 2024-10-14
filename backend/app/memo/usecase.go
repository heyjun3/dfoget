package memo

import (
	"context"

	domain "github.com/heyjun3/dforget/backend/domain/memo"
	"github.com/heyjun3/dforget/backend/lib"
)

type MemoRepositoryInterface interface {
	Find(context.Context) ([]*domain.Memo, error)
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
