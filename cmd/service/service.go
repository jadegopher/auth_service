package service

import (
	"auth/cmd/db"
	"auth/proto"
	"go.uber.org/zap"
)

const (
	rsaPrivateKey = "RSA PRIVATE KEY"
)

type handler struct {
	logger     *zap.Logger
	usersDB    db.IUsers
	sessionsDB db.ISessions
}

func New(logger *zap.Logger, users db.IUsers, session db.ISessions) proto.AuthServiceServer {
	return &handler{
		logger:     logger,
		usersDB:    users,
		sessionsDB: session,
	}
}
