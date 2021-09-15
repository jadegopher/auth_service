package service

import (
	"go.uber.org/zap"

	"auth/internal/core/ports"
	"auth/internal/core/service/encription"
	"auth/internal/core/service/encription/rsa"
	"auth/internal/core/service/token"
)

type TokenCreator interface {
	Create(payload map[string]interface{}, key interface{}) (_ string, err error)
}

type KeyGenerator interface {
	Generate() (key encription.IKey, err error)
}

type Service struct {
	logger       *zap.Logger
	usersDB      ports.IUsers
	sessionsDB   ports.ISessions
	tokenCreator TokenCreator
	keyGenerator KeyGenerator
}

func New(logger *zap.Logger, users ports.IUsers, session ports.ISessions) *Service {
	return &Service{
		logger:       logger,
		usersDB:      users,
		sessionsDB:   session,
		tokenCreator: token.New(logger),
		keyGenerator: rsa.NewGenerator(logger),
	}
}
