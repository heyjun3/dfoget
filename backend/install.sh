#!/bin/bash

go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

# install atlas
curl -sSf https://atlasgo.sh | sh

# install wire
go install github.com/google/wire/cmd/wire@latest