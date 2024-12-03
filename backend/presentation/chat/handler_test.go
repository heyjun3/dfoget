package chat_test

import (
	"context"
	// "database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"

	chatv1 "github.com/heyjun3/dforget/backend/gen/api/chat/v1"
	"github.com/heyjun3/dforget/backend/gen/api/chat/v1/chatv1connect"
	model "github.com/heyjun3/dforget/backend/infra/chat"
	"github.com/heyjun3/dforget/backend/lib"
	"github.com/heyjun3/dforget/backend/test"

	"github.com/heyjun3/dforget/backend/presentation"
	// "github.com/heyjun3/dforget/backend/presentation/chat"
	"github.com/heyjun3/dforget/backend/server"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	// "github.com/uptrace/bun/dialect/pgdialect"
	// "github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func newTestServer() (*httptest.Server, *bun.DB, func()) {
	conf := server.NewConfig(server.WithDBName("test"), server.WithPubKey(test.PublicKey))
	bundb := server.InitDBConn(conf)
	mux := presentation.NewServer(conf)
	// sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.DBDSN())))
	// bundb := bun.NewDB(sqldb, pgdialect.New())
	// chat := chat.InitChatHandler(bundb)
	// mux := http.NewServeMux()
	// path, handler := chatv1connect.NewChatServiceHandler(chat)
	// mux.Handle(path, handler)
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

	t.Run("create room", func(t *testing.T) {
		userID, _ := uuid.NewV7()
		ctx := lib.SetSubKey(context.Background(), userID.String())
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
