package chat

import (
	"context"
	"fmt"
	"log/slog"
)

type CreateRoomRepositoryInterface interface {
	Exists(context.Context, string) (bool, error)
	Save(context.Context, *Room) error
}
type CreateRoomService struct {
	createRoomRepository CreateRoomRepositoryInterface
}

func NewCreateRoomService(createRoomRepository CreateRoomRepositoryInterface) *CreateRoomService {
	return &CreateRoomService{
		createRoomRepository: createRoomRepository,
	}
}

func (s *CreateRoomService) Execute(ctx context.Context, name string) (
	*Room, error,
) {
	exists, err := s.createRoomRepository.Exists(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("existing room name")
	}
	room, err := NewRoom(name)
	if err != nil {
		return nil, err
	}
	if err := s.createRoomRepository.Save(ctx, room); err != nil {
		slog.WarnContext(ctx, err.Error())
		return nil, fmt.Errorf("failed to save room")
	}
	return room, err
}
