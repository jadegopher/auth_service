package service

import (
	"auth/internal/entities"
	"auth/internal/infrastructure/db"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
)

// CheckToken TODO: unit test
func (s *service) CheckToken(ctx context.Context, token string) (err error) {
	var userID int64
	if userID, err = s.sessionsDB.GetUserIDByToken(token); err != nil {
		if err == db.NotFound {
			s.logger.Info("Token expired", zap.String("token", token))
			return ErrTokenExpired
		}
		s.logger.Error("error sessionsDB.GetUserIDByToken", zap.Error(err))
		return err
	}

	var userInfo *entities.User
	if userInfo, err = s.usersDB.SelectByUserID(ctx, userID); err != nil {
		s.logger.Error("error usersDB.SelectByUserID", zap.Error(err))
		if err == db.NotFound {
			return ErrAccountNotFound
		}
		return err
	}

	var tokenUserID int64
	if tokenUserID, err = s.getUserIDFromToken(userInfo, []byte(token)); err != nil {
		return err
	}

	if tokenUserID != userID {
		return ErrInvalidToken
	}

	return nil
}

func (s *service) getUserIDFromToken(userInfo *entities.User, payload []byte) (_ int64, err error) {
	var block *pem.Block
	block, _ = pem.Decode([]byte(userInfo.Password))
	if block == nil || block.Type != rsaPrivateKey {
		s.logger.Error("error pem.Decode: wrong block type")
		return 0, err
	}

	var privateKey *rsa.PrivateKey
	if privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		s.logger.Error("error x509.ParsePKCS1PrivateKey", zap.Error(err))
		return 0, err
	}

	var token jwt.Token
	if token, err = jwt.Parse(
		payload,
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, &privateKey.PublicKey),
	); err != nil {
		s.logger.Error("error jwt.Parse", zap.Error(err))
		return 0, ErrInvalidToken
	}

	var (
		value interface{}
		ok    bool
	)
	if value, ok = token.Get(entities.UserIDField); !ok {
		s.logger.Error("error token.Get UserIDField", zap.Error(err))
		return 0, ErrInvalidToken
	}

	var cast float64
	if cast, ok = value.(float64); !ok {
		s.logger.Error("error cast UserIDField to float64", zap.Error(err))
		return 0, ErrInvalidToken
	}

	return int64(cast), nil
}
