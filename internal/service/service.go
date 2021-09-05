package service

import (
	"auth/internal/entities"
	"go.uber.org/zap"
)

const (
	rsaPrivateKey = "RSA PRIVATE KEY"
)

type service struct {
	logger     *zap.Logger
	usersDB    entities.IUsers
	sessionsDB entities.ISessions
}

func New(logger *zap.Logger, users entities.IUsers, session entities.ISessions) *service {
	return &service{
		logger:     logger,
		usersDB:    users,
		sessionsDB: session,
	}
}
