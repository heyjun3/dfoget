package chat

import (
	"context"

	"github.com/google/uuid"
	"github.com/heyjun3/dforget/backend/domain/chat"
)

type RoomRepository interface {
	DeleteById(context.Context, uuid.UUID) error
	GetRoom(context.Context, uuid.UUID) (*chat.Room, error)
	GetRoomsWithoutMessage(context.Context) ([]*chat.RoomWithoutMessage, error)
	Save(context.Context, *chat.Room) error
}

type RoomUsecase struct {
	createRoomService chat.CreateRoomService
	roomRepository    RoomRepository
}

func NewRoomUsecase(
	createRoomService *chat.CreateRoomService,
	roomRepository RoomRepository,
) *RoomUsecase {
	return &RoomUsecase{
		createRoomService: *createRoomService,
		roomRepository:    roomRepository,
	}
}

func (u *RoomUsecase) CreateRoom(ctx context.Context, name string) (*chat.Room, error) {
	room, err := u.createRoomService.Execute(ctx, name)
	if err != nil {
		return nil, err
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

func (u *RoomUsecase) AddMessage(ctx context.Context, roomId uuid.UUID, text string) error {
	room, err := u.roomRepository.GetRoom(ctx, roomId)
	if err != nil {
		return err
	}
	room.AddMessage(ctx, text)
	return u.roomRepository.Save(ctx, room)
}
