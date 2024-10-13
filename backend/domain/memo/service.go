package memo

import (
	"context"

	"github.com/google/uuid"
)

type MemoRepositoryInterface interface {
	GetById(context.Context, uuid.UUID) (*Memo, error)
	Save(context.Context, []*Memo) ([]*Memo, error)
}

type RegisterMemoService struct {
	memoRepository MemoRepositoryInterface
}

func NewRegisterMemoService(memoRepository MemoRepositoryInterface) *RegisterMemoService {
	return &RegisterMemoService{
		memoRepository: memoRepository,
	}
}

func (s RegisterMemoService) Execute(
	ctx context.Context, sub uuid.UUID, id *string, title, text string) (
	*Memo, error,
) {
	var opts []Option
	var memo *Memo
	if id != nil {
		uu, err := uuid.Parse(*id)
		if err != nil {
			return nil, err
		}
		memo, err = s.memoRepository.GetById(ctx, uu)
		if err != nil {
			return nil, err
		}
		memo.Title = title
		memo.Text = text
	} else {
		memoPtr, err := NewMemo(title, text, sub, opts...)
		if err != nil {
			return nil, err
		}
		memo = memoPtr
	}
	_, err := s.memoRepository.Save(ctx, []*Memo{memo})
	if err != nil {
		return nil, err
	}
	return memo, nil
}
