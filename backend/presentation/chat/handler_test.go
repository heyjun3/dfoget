package chat_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	cfg "github.com/heyjun3/dforget/backend/config"
	chatv1 "github.com/heyjun3/dforget/backend/gen/api/chat/v1"
	"github.com/heyjun3/dforget/backend/gen/api/chat/v1/chatv1connect"
	model "github.com/heyjun3/dforget/backend/infra/chat"
	"github.com/heyjun3/dforget/backend/presentation"
	"github.com/heyjun3/dforget/backend/server"
	"github.com/heyjun3/dforget/backend/test"
)

func newTestServer() (*httptest.Server, *bun.DB, func()) {
	conf := cfg.NewConfig(cfg.WithDBName("test"), cfg.WithPubKey(test.PublicKey))
	bundb := server.InitDBConn(conf)
	mux := presentation.NewServer(conf)
	srv := httptest.NewServer(h2c.NewHandler(mux, &http2.Server{}))
	return srv, bundb, func() {
		srv.Close()
		bundb.Close()
	}
}

func TestChatHandler(t *testing.T) {
	srv, db, tireDown := newTestServer()
	if _, err := db.NewDelete().Model(&model.RoomDM{}).Where("1 = 1").Exec(context.Background()); err != nil {
		panic(err)
	}
	defer tireDown()
	client := chatv1connect.NewChatServiceClient(
		http.DefaultClient,
		srv.URL,
		connect.WithInterceptors(test.NewSetCookieInterceptor()),
	)

	t.Run("create room and get room", func(t *testing.T) {
		ctx := context.Background()
		createRes, err := client.CreateRoom(
			ctx,
			connect.NewRequest(&chatv1.CreateRoomRequest{
				Name: "test room",
			}),
		)
		assert.NoError(t, err)
		fmt.Println("res", createRes.Msg.Room)
		assert.Equal(t, "test room", createRes.Msg.Room.Name)

		res, err := client.GetRooms(
			ctx,
			connect.NewRequest(&chatv1.GetRoomsRequest{}),
		)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(res.Msg.Rooms))
		assert.Equal(t, "test room", res.Msg.Rooms[0].Name)

		getRoomRes, err := client.GetRoom(
			ctx,
			connect.NewRequest(&chatv1.GetRoomRequest{
				Id: createRes.Msg.Room.Id,
			}),
		)

		assert.NoError(t, err)
		assert.Equal(t, "test room", getRoomRes.Msg.GetRoom().Name)
		assert.Equal(t, 0, len(getRoomRes.Msg.GetMessages()))

		sendRes, err := client.SendMessage(ctx, connect.NewRequest(
			&chatv1.SendMessageRequest{
				RoomId: createRes.Msg.Room.Id,
				Text:   "test text",
			},
		))

		assert.NoError(t, err)
		assert.Equal(t, "test text", sendRes.Msg.Message.Text)
	})
}
