package service

import (
	"go.uber.org/zap"

	"auth/internal/entities"
	"auth/internal/service/key"
	"auth/internal/service/key/aes"
	"auth/internal/service/key/rsa"
	"auth/internal/service/token

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

type Encryptor interface {
	Encrypt(key []byte, plaintext string) (string, error)
}

type Service struct {
	logger       *zap.Logger
	usersDB      ports.IUsers
	sessionsDB   ports.ISessions
	cypherKey    []byte
	tokenCreator TokenCreator
	keyGenerator KeyGenerator
	encryptor    Encryptor
}

func New(logger *zap.Logger, cypherKey []byte, users ports.IUsers, session ports.ISessions) *Service {
	return &Service{
		logger:       logger,
		cypherKey:    cypherKey,
		usersDB:      users,
		sessionsDB:   session,
		tokenCreator: token.New(logger),
		keyGenerator: rsa.NewGenerator(logger),
		encryptor:    aes.NewEncrypt(logger),
	}
}
