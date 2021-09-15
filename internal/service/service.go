package service

import (
	"auth/internal/entities"
	"auth/internal/service/key"
	"auth/internal/service/key/aes"
	"auth/internal/service/key/rsa"
	"auth/internal/service/token"
	"go.uber.org/zap"
)

type TokenCreator interface {
	Create(payload map[string]interface{}, key interface{}) (_ string, err error)
}

type KeyGenerator interface {
	Generate() (key key.IKey, err error)
}

type AESEncryptor interface {
	AESEncrypt(key []byte, plaintext string) (string, error)
}

type service struct {
	logger       *zap.Logger
	usersDB      entities.IUsers
	sessionsDB   entities.ISessions
	tokenCreator TokenCreator
	keyGenerator KeyGenerator
	AESEncryptor AESEncryptor
}

func New(logger *zap.Logger, users entities.IUsers, session entities.ISessions) *service {
	return &service{
		logger:       logger,
		usersDB:      users,
		sessionsDB:   session,
		tokenCreator: token.New(logger),
		keyGenerator: rsa.NewGenerator(logger),
		AESEncryptor: aes.NewEncrypt(logger),
	}
}
