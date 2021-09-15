package main

import (
	_ "github.com/golang/mock/mockgen/model"
	_ "google.golang.org/grpc"
)

//go:generate protoc -I. -I$GOPATH\src -I$GOPATH\src\github.com\grpc-ecosystem\grpc-gateway\third_party\googleapis --go_out=plugins=grpc,paths=source_relative:./ --grpc-gateway_out=logtostderr=true:./ --swagger_out=allow_merge=true,merge_file_name=api:./proto ./proto/*.proto
//go:generate mockgen --build_flags=--mod=mod -destination=mocks/sessions.go -package=mocks auth/internal/core/ports ISessions
//go:generate mockgen --build_flags=--mod=mod -destination=mocks/users.go -package=mocks auth/internal/core/ports IUsers
//go:generate mockgen --build_flags=--mod=mod -destination=mocks/token_creator.go -package=mocks auth/internal/core/service TokenCreator
//go:generate mockgen --build_flags=--mod=mod -destination=mocks/key_generator.go -package=mocks auth/internal/core/service KeyGenerator
