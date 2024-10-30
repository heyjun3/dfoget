package chat

import (
	"context"
	"fmt"
	"log/slog"
)

type CreateRoomRepositoryInterface interface {
	Exists(context.Context, string) (bool, error)
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
	if exists, err := s.createRoomRepository.Exists(ctx, name); exists || err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, fmt.Errorf("existing room name")
	}
	room, err := NewRoom(name)
	if err != nil {
		return nil, err
	}
	return room, err
}
