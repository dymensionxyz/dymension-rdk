//go:build tools
// +build tools

package tools

import (
	_ "github.com/cosmos/gogoproto/gogoproto"
	_ "github.com/cosmos/gogoproto/jsonpb"
	_ "github.com/cosmos/gogoproto/proto"
	_ "github.com/cosmos/gogoproto/protoc-gen-gogo"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
