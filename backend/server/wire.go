// go:build wireinject
//go:build wireinject
// +build wireinject

package server

import (
	"database/sql"
	"net/http"

	"github.com/google/wire"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	chatapp "github.com/heyjun3/dforget/backend/app/chat"
	memoapp "github.com/heyjun3/dforget/backend/app/memo"
	chatdomain "github.com/heyjun3/dforget/backend/domain/chat"
	"github.com/heyjun3/dforget/backend/domain/memo"
	chatinfra "github.com/heyjun3/dforget/backend/infra/chat"
	"github.com/heyjun3/dforget/backend/presentation/chat"
)

func initializeMemoHandler(db *bun.DB) *MemoHandler {
	wire.Build(
		NewMemoRepository,
		memo.NewRegisterMemoService,
		memoapp.NewMemoUsecase,
		NewMemoHandler,
		wire.Bind(new(memo.MemoRepositoryInterface), new(*MemoRepository)),
		wire.Bind(new(memoapp.MemoRepositoryInterface), new(*MemoRepository)),
	)
	return nil
}

func initializeChatHandler(db *bun.DB) *chat.ChatServiceHandler {
	wire.Build(
		chat.NewChatServiceHandler,
		chatapp.NewRoomUsecase,
		chatdomain.NewCreateRoomService,
		chatinfra.NewChatRepository,
		wire.Bind(new(chatdomain.CreateRoomRepositoryInterface), new(*chatinfra.ChatRepository)),
		wire.Bind(new(chatapp.RoomRepository), new(*chatinfra.ChatRepository)),
	)
	return nil
}

func initializeOIDCHandler(conf Config) *OIDCHandler {
	wire.Build(
		provideHttpClient,
		NewOIDCHandler,
	)
	return nil
}

func provideHttpClient() httpClient {
	return &http.Client{}
}

type DBConfigIF interface {
	DBDSN() string
}

func provideOpenDB(conf DBConfigIF) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.DBDSN())))
	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}
func InitDBConn(conf DBConfigIF) *bun.DB {
	wire.Build(
		provideOpenDB,
	)
	return nil
}
