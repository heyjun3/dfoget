//go:build wireinject
// +build wireinject

package chat

import (
	"github.com/google/wire"
	"github.com/uptrace/bun"

	chatapp "github.com/heyjun3/dforget/backend/app/chat"
	chatinfra "github.com/heyjun3/dforget/backend/infra/chat"
)

func InitChatHandler(db *bun.DB) *ChatServiceHandler {
	wire.Build(
		NewChatServiceHandler,
		chatapp.NewRoomUsecase,
		chatinfra.NewChatRepository,
		wire.Bind(new(chatapp.RoomRepository), new(*chatinfra.ChatRepository)),
	)
	return nil
}
