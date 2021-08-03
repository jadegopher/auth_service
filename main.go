package main

import (
	"auth/cmd/db/postgresql"
	"auth/cmd/db/postgresql/users"
	"auth/cmd/db/redis"
	"auth/cmd/db/redis/sessions"
	"auth/cmd/model"
	"auth/cmd/service"
	"auth/proto"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	pConn, err := postgresql.NewConnection(model.Database{
		IP:       "127.0.0.1",
		Port:     "5432",
		Name:     "postgres",
		User:     "postgres",
		Password: "",
	})
	if err != nil {
		panic(err)
	}

	rConn, err := redis.NewConnection(model.Database{
		IP:   "0.0.0.0",
		Port: "6379",
	})
	if err != nil {
		panic(err)
	}

	micro := service.New(logger, users.New(pConn), sessions.New(rConn))

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	mux := runtime.NewServeMux()

	if err = proto.RegisterAuthServiceHandlerServer(ctx, mux, micro); err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("HTTP API Listening on %s port", "3000"))

	if err = http.ListenAndServe(fmt.Sprintf(":%s", "3000"), mux); err != nil {
		panic(err)
	}
}
