package memo

import (
	"context"

	domain "github.com/heyjun3/dforget/backend/domain/memo"
	"github.com/heyjun3/dforget/backend/lib"
)

func NewMemoUsecase(registerMemoService *domain.RegisterMemoService) *MemoUsecase {
	return &MemoUsecase{
		registerMemoService: registerMemoService,
	}
}

type MemoUsecase struct {
	registerMemoService *domain.RegisterMemoService
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
