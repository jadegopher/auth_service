package aes

import (
	"crypto/aes"
	"encoding/hex"
	"go.uber.org/zap"
)

type encrypt struct {
	logger *zap.Logger
}

func NewEncrypt(logger *zap.Logger) *encrypt {
	return &encrypt{logger: logger}
}

// AESEncrypt Encrypt text using AES algorithm
func (e *encrypt) AESEncrypt(key []byte, plaintext string) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		e.logger.Error("error aes.NewCipher", zap.Error(err))
		return "", err
	}

	out := make([]byte, len(plaintext))
	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out), nil
}
