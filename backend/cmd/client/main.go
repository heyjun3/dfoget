package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"

	memov1 "github.com/heyjun3/dforget/backend/gen/api/memo/v1"
	"github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
)

func main() {
	client := memov1connect.NewMemoServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPCWeb(),
	)
	res, err := client.GetMemo(
		context.Background(),
		connect.NewRequest(&memov1.GetMemoRequest{}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg)
}
