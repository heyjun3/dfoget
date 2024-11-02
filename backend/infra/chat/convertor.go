package chat

import (
	"github.com/heyjun3/dforget/backend/domain/chat"
)

func dmToRoom(dm *RoomDM) *chat.Room {
	r, _ := chat.NewRoom(dm.Name, chat.WithReconstructRoom(dm.ID, dmToMessages(dm.Messages), dm.CreatedAt))
	return r
}

func dmToRoomWithoutMessage(dm *RoomDM) *chat.RoomWithoutMessage {
	return &chat.RoomWithoutMessage{
		ID:        dm.ID,
		Name:      dm.Name,
		CreatedAt: dm.CreatedAt,
	}
}
func dmToRoomWithoutMessages(dm []*RoomDM) []*chat.RoomWithoutMessage {
	rooms := make([]*chat.RoomWithoutMessage, 0, len(dm))
	for _, r := range dm {
		rooms = append(rooms, dmToRoomWithoutMessage(r))
	}
	return rooms
}

func dmToMessage(dm *MessageDM) *chat.Message {
	m, _ := chat.NewMessage(
		dm.UserID,
		dm.RoomID,
		dm.Text,
		chat.WithReconstruct(dm.ID, dm.CreatedAt),
	)
	return m
}

func dmToMessages(dm []*MessageDM) []chat.Message {
	messages := make([]chat.Message, 0, len(dm))
	for _, d := range dm {
		messages = append(messages, *dmToMessage(d))
	}
	return messages
}

func roomToDM(room *chat.Room) *RoomDM {
	id, name, messages, createdAt := room.Get()
	dm := &RoomDM{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		Messages:  messagesToDM(messages),
	}
	return dm
}

func messageToDM(message *chat.Message) *MessageDM {
	id, userID, roomID, text, createdAt := message.Get()
	return &MessageDM{
		ID:        id,
		UserID:    userID,
		RoomID:    roomID,
		Text:      text,
		CreatedAt: createdAt,
	}
}
func messagesToDM(messages []chat.Message) []*MessageDM {
	dm := make([]*MessageDM, 0, len(messages))
	for _, message := range messages {
		dm = append(dm, messageToDM(&message))
	}
	return dm
}
