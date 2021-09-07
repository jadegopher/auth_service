package service

import (
	"auth/internal/entities"
	"auth/internal/service/key"
	"auth/internal/service/token"
	"go.uber.org/zap"
)

type TokenCreator interface {
	Create(payload map[string]interface{}, key interface{}) (_ string, err error)
}

type KeyGenerator interface {
	Generate() (key key.IKey, err error)
}

type IKey interface {
	String() string
	Interface() interface{}
}

type service struct {
	logger       *zap.Logger
	usersDB      entities.IUsers
	sessionsDB   entities.ISessions
	tokenCreator TokenCreator
	keyGenerator KeyGenerator
}

func New(logger *zap.Logger, users entities.IUsers, session entities.ISessions) *service {
	return &service{
		logger:       logger,
		usersDB:      users,
		sessionsDB:   session,
		tokenCreator: token.New(logger),
		keyGenerator: key.New(logger),
	}
}
