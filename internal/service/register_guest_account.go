package service

import (
	"auth/internal/entities"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
	"time"
)

// RegisterGuestAccount TODO: unit test
func (s *service) RegisterGuestAccount(ctx context.Context, username string) (_ string, err error) {
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		s.logger.Error("error rsa.GenerateKey", zap.Error(err))
		return "", err
	}

	var userID int64
	// TODO add transaction and if some error rollback it
	if userID, err = s.usersDB.Insert(
		ctx,
		username,
		string(pem.EncodeToMemory(
			&pem.Block{
				Type:  rsaPrivateKey,
				Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
			},
		)),
	); err != nil {
		s.logger.Error("error usersDB.Insert", zap.Error(err))
		return "", err
	}

	var token = jwt.New()
	if err = token.Set(jwt.IssuedAtKey, time.Now().UTC().Unix()); err != nil {
		s.logger.Error("error set IssuedAtKey", zap.Error(err))
		return "", err
	}

	if err = token.Set(entities.UserIDField, userID); err != nil {
		s.logger.Error("error set UserIDField", zap.Error(err))
		return "", err
	}

	// Sign the token and generate a payload
	var signedToken []byte
	if signedToken, err = jwt.Sign(token, jwa.RS256, privateKey); err != nil {
		s.logger.Error("error jwt.Sign", zap.Error(err))
		return "", err
	}

	if err = s.sessionsDB.Insert(string(signedToken), userID, time.Hour); err != nil {
		s.logger.Error("error sessionsDB.Insert", zap.Error(err))
		return "", err
	}

	return string(signedToken), nil
}
