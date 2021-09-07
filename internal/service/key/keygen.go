package key

import (
	"crypto/rand"
	"crypto/rsa"
	"go.uber.org/zap"
)

type keygen struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *keygen {
	return &keygen{logger: logger}
}

// Generate generates rsa key pair
// TODO think about interface service.IKey
func (k *keygen) Generate() (key IKey, err error) {
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		k.logger.Error("error rsa.GenerateKey", zap.Error(err))
		return nil, err
	}

	return &Key{privateKey: privateKey}, nil
}
