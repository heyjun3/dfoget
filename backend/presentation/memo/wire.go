//go:build wireinject

package memo

import (
	"github.com/google/wire"
	"github.com/uptrace/bun"

	memoapp "github.com/heyjun3/dforget/backend/app/memo"
	"github.com/heyjun3/dforget/backend/domain/memo"
	memodm "github.com/heyjun3/dforget/backend/infra/memo"
)

func InitializeMemoHandler(db *bun.DB) *MemoHandler {
	wire.Build(
		memodm.NewMemoRepository,
		memo.NewRegisterMemoService,
		memoapp.NewMemoUsecase,
		NewMemoHandler,
		wire.Bind(new(memo.MemoRepositoryInterface), new(*memodm.MemoRepository)),
		wire.Bind(new(memoapp.MemoRepositoryInterface), new(*memodm.MemoRepository)),
	)
	return nil
}
