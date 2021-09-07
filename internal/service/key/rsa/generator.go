package rsa

import (
	"auth/internal/service/key"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"go.uber.org/zap"
)

const rsaPrivateKey = "RSA PRIVATE KEY"

type keyContainer struct {
	privateKey *rsa.PrivateKey
}

func (k *keyContainer) String() string {
	return string(pem.EncodeToMemory(
		&pem.Block{
			Type:  rsaPrivateKey,
			Bytes: x509.MarshalPKCS1PrivateKey(k.privateKey),
		},
	))
}

func (k *keyContainer) Interface() interface{} {
	return k.privateKey
}

type keygen struct {
	logger *zap.Logger
}

func NewGenerator(logger *zap.Logger) *keygen {
	return &keygen{logger: logger}
}

// Generate generates rsa key pair
func (k *keygen) Generate() (key key.IKey, err error) {
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		k.logger.Error("error rsa.GenerateKey", zap.Error(err))
		return nil, err
	}

	return &keyContainer{privateKey: privateKey}, nil
}
