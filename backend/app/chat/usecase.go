package chat

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/heyjun3/dforget/backend/domain/chat"
)

type RoomRepository interface {
	Save(context.Context, *chat.Room) error
	GetRoom(context.Context, uuid.UUID) (*chat.Room, error)
	GetRoomsWithoutMessage(context.Context) ([]*chat.RoomWithoutMessage, error)
	Exists(context.Context, string) (bool, error)
	DeleteById(context.Context, uuid.UUID) error
}

type RoomUsecase struct {
	roomRepository RoomRepository
}

func NewRoomUsecase(
	roomRepository RoomRepository,
) *RoomUsecase {
	return &RoomUsecase{
		roomRepository: roomRepository,
	}
}

func (u *RoomUsecase) CreateRoom(ctx context.Context, name string) (*chat.Room, error) {
	exists, err := u.roomRepository.Exists(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("existing room name")
	}
	room, err := chat.NewRoom(name)
	if err != nil {
		return nil, err
	}
	if err := u.roomRepository.Save(ctx, room); err != nil {
		slog.WarnContext(ctx, err.Error())
		return nil, fmt.Errorf("failed to save room")
	}
	return room, err
}

func (u *RoomUsecase) FetchRoom(ctx context.Context, id uuid.UUID) (*chat.Room, error) {
	return u.roomRepository.GetRoom(ctx, id)
}

func (u *RoomUsecase) FetchRooms(ctx context.Context) ([]*chat.RoomWithoutMessage, error) {
	return u.roomRepository.GetRoomsWithoutMessage(ctx)
}

func (u *RoomUsecase) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	if err := u.roomRepository.DeleteById(ctx, id); err != nil {
		return err
	}
	return nil
}

func (u *RoomUsecase) AddMessage(ctx context.Context, roomId uuid.UUID, text string) (*chat.Message, error) {
	room, err := u.roomRepository.GetRoom(ctx, roomId)
	if err != nil {
		return nil, err
	}
	message, err := room.AddMessage(ctx, text)
	if err != nil {
		return nil, err
	}
	if err := u.roomRepository.Save(ctx, room); err != nil {
		return nil, err
	}
	return message, nil
}
