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

type Encryptor interface {
	Encrypt(key []byte, plaintext string) (string, error)
}

type service struct {
	logger       *zap.Logger
	cypherKey    []byte
	usersDB      entities.IUsers
	sessionsDB   entities.ISessions
	tokenCreator TokenCreator
	keyGenerator KeyGenerator
	encryptor    Encryptor
}

func New(logger *zap.Logger, cypherKey []byte, users entities.IUsers, session entities.ISessions) *service {
	return &service{
		logger:       logger,
		cypherKey:    cypherKey,
		usersDB:      users,
		sessionsDB:   session,
		tokenCreator: token.New(logger),
		keyGenerator: rsa.NewGenerator(logger),
		encryptor:    aes.NewEncrypt(logger),
	}
}
