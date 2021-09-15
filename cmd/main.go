package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"

	"auth/internal/adapters/handlers"
	"auth/internal/core/entities"
	"auth/internal/core/service"
	"auth/internal/infrastructure/db/postgresql"
	"auth/internal/infrastructure/db/postgresql/users"
	rDB "auth/internal/infrastructure/db/redis"
	"auth/internal/infrastructure/db/redis/sessions"
	"auth/proto"
)

func main() {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ctx, cancelFunc := context.WithCancel(context.Background())

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	pConn, err := postgresql.NewConnection(entities.Database{
		IP:       "127.0.0.1",
		Port:     "5432",
		Name:     "postgres",
		User:     "postgres",
		Password: "postgres",
	})
	if err != nil {
		logger.Error("Error postgresql.NewConnection", zap.Error(err))
		panic(err)
	}
	defer closeConn(logger, pConn, "postgres")

	rConn, err := rDB.NewConnection(entities.Database{
		IP:   "0.0.0.0",
		Port: "6379",
	})
	if err != nil {
		logger.Error("Error redis.NewConnection", zap.Error(err))
		panic(err)
	}
	defer closeConn(logger, rConn, "redis")

	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%s", "3001"),
	}
	defer func() {
		if err = httpServer.Shutdown(context.Background()); err != nil {
			logger.Error("httpServer.Shutdown", zap.Error(err))
		}
		logger.Info("Server exited properly")
	}()

	go func() {
		_ = run(ctx, logger, pConn, rConn, httpServer)
	}()

	for range termChan {
		cancelFunc()
	}

	logger.Info("Starting graceful shutdown")
}

func run(
	ctx context.Context,
	logger *zap.Logger,
	pConn *sql.DB,
	rConn *redis.Client,
	server *http.Server,
) (err error) {
	authService := service.New(logger, users.New(pConn), sessions.New(rConn))
	handler := handlers.New(authService)
	mux := runtime.NewServeMux()
	server.Handler = mux

	if err = proto.RegisterAuthServiceHandlerServer(ctx, mux, handler); err != nil {
		logger.Error("Error proto.RegisterAuthServiceHandlerServer", zap.Error(err))
		return err
	}

	logger.Info(fmt.Sprintf("Server started on %s", server.Addr))

	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Error server.ListenAndServe", zap.Error(err))
		return err
	}

	return nil
}

func closeConn(logger *zap.Logger, closer io.Closer, name string) {
	if err := closer.Close(); err != nil {
		logger.Error("Fatal closing connection", zap.Error(err))
	}
	logger.Info("Connection closed", zap.String("source_name", name))
}
