package chat_handler

import (
	"context"

	"connectrpc.com/connect"
	app "github.com/heyjun3/dforget/backend/app/chat"
	chatv1 "github.com/heyjun3/dforget/backend/gen/api/chat/v1"
	"github.com/heyjun3/dforget/backend/gen/api/chat/v1/chatv1connect"
)

var _ chatv1connect.ChatServiceHandler = (*ChatServiceHandler)(nil)

type ChatServiceHandler struct {
	roomUsecase app.RoomUsecase
}

func NewChatServiceHandler(roomUsecase app.RoomUsecase) *ChatServiceHandler {
	return &ChatServiceHandler{
		roomUsecase: roomUsecase,
	}
}

func (c *ChatServiceHandler) GetRooms(ctx context.Context, req *connect.Request[chatv1.GetRoomsRequest]) (
	*connect.Response[chatv1.GetRoomsResponse], error,
) {
	rooms, err := c.roomUsecase.FetchRooms(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, nil)
	}
	roomsDto := make([]*chatv1.Room, len(rooms))
	for i := range rooms {
		roomsDto[i] = &chatv1.Room{
			Id:   rooms[i].ID.String(),
			Name: rooms[i].Name,
		}
	}
	return connect.NewResponse(
		&chatv1.GetRoomsResponse{
			Rooms: roomsDto,
		},
	), nil
}

func (c *ChatServiceHandler) CreateRoom(ctx context.Context, req *connect.Request[chatv1.CreateRoomRequest]) (
	*connect.Response[chatv1.CreateRoomResponse], error,
) {
	return nil, nil
}
