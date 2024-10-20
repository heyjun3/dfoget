package main

import (
	"database/sql"
	"os"

	"github.com/heyjun3/dforget/backend/infra/chat"
	"github.com/heyjun3/dforget/backend/server"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func ptr[T any](t T) *T {
	return &t
}

func main() {
	dsn := "postgres://dev:dev@postgres:5432/dforget?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	models := []server.Models{
		{Model: &server.MemoDM{}, Fkey: nil},
		{Model: &chat.RoomDM{}, Fkey: nil},
		{Model: &chat.MessageDM{}, Fkey: ptr(`("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE`)},
	}
	query := server.ModelsToBytes(db, models)
	f, err := os.Create("./schema.sql")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(query)
}
