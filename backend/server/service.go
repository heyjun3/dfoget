package server

import (
	"context"

	"github.com/google/uuid"
)

type RegisterMemoService struct {
	memoRepository *MemoRepository
}

func NewRegisterMemoService(memoRepository MemoRepository) *RegisterMemoService {
	return &RegisterMemoService{
		memoRepository: &memoRepository,
	}
}

func (s RegisterMemoService) execute(
	ctx context.Context, sub uuid.UUID, id *string, title, text string) (
	*Memo, error,
) {
	var opts []Option
	var memo Memo
	if id != nil {
		uu, err := uuid.Parse(*id)
		if err != nil {
			return nil, err
		}
		memo, err = s.memoRepository.GetById(ctx, sub, uu)
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
		memo = *memoPtr
	}
	_, err := s.memoRepository.Save(context.Background(), []Memo{memo})
	if err != nil {
		return nil, err
	}
	return &memo, nil
}
