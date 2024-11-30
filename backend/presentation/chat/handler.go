package chat

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/google/uuid"

	app "github.com/heyjun3/dforget/backend/app/chat"
	chatv1 "github.com/heyjun3/dforget/backend/gen/api/chat/v1"
	"github.com/heyjun3/dforget/backend/gen/api/chat/v1/chatv1connect"
	"github.com/heyjun3/dforget/backend/lib"
)

var _ chatv1connect.ChatServiceHandler = (*ChatServiceHandler)(nil)

type ChatServiceHandler struct {
	roomUsecase app.RoomUsecase
}

func NewChatServiceHandler(roomUsecase *app.RoomUsecase) *ChatServiceHandler {
	return &ChatServiceHandler{
		roomUsecase: *roomUsecase,
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

func (c *ChatServiceHandler) GetRoom(ctx context.Context, req *connect.Request[chatv1.GetRoomRequest]) (
	*connect.Response[chatv1.GetRoomResponse], error,
) {
	id, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		slog.InfoContext(ctx, "failed parse id")
		return nil, connect.NewError(connect.CodeInternal, nil)
	}
	room, err := c.roomUsecase.FetchRoom(ctx, id)
	if err != nil {
		slog.InfoContext(ctx, "failed fetch room", "error", err)
		return nil, connect.NewError(connect.CodeInternal, nil)
	}
	roomID, name, messages, _ := room.Get()
	messageDTO := make([]*chatv1.Message, 0, len(messages))
	for _, msg := range messages {
		msgID, userID, _, text, _ := msg.Get()
		messageDTO = append(messageDTO, &chatv1.Message{
			Id:     msgID.String(),
			UserId: lib.Ptr(userID.String()),
			Text:   text,
		})
	}
	return connect.NewResponse(
		&chatv1.GetRoomResponse{
			Room: &chatv1.Room{
				Id:   roomID.String(),
				Name: name,
			},
			Messages: messageDTO,
		},
	), nil
}

func (c *ChatServiceHandler) CreateRoom(ctx context.Context, req *connect.Request[chatv1.CreateRoomRequest]) (
	*connect.Response[chatv1.CreateRoomResponse], error,
) {
	room, err := c.roomUsecase.CreateRoom(ctx, req.Msg.Name)
	if err != nil {
		slog.InfoContext(ctx, err.Error())
		return nil, connect.NewError(connect.CodeInternal, nil)
	}
	id, name, _, _ := room.Get()
	return connect.NewResponse(
		&chatv1.CreateRoomResponse{
			Room: &chatv1.Room{
				Id:   id.String(),
				Name: name,
			},
		},
	), nil
}

func (c *ChatServiceHandler) SendMessage(ctx context.Context, req *connect.Request[chatv1.SendMessageRequest]) (
	*connect.Response[chatv1.SendMessageResponse], error,
) {
	return nil, nil
}
