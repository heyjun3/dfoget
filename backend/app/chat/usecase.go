package chat

import (
	"context"

	"github.com/heyjun3/dforget/backend/domain/chat"
)

type RoomUsecase struct {
	createRoomService chat.CreateRoomService
}

func (u *RoomUsecase) CreateRoom(ctx context.Context, name string) (*chat.Room, error) {
	room, err := u.createRoomService.Execute(ctx, name)
	if err != nil {
		return nil, err
	}
	return room, err
}
