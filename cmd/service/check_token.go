package service

import (
	"auth/cmd/db"
	"auth/proto"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
)

func (h *handler) CheckToken(
	ctx context.Context,
	req *proto.CheckTokenRequest,
) (resp *proto.CheckTokenResponse, err error) {
	var userID int64
	if userID, err = h.sessionsDB.GetUserIDByToken(req.Token); err != nil {
		if err == db.NotFound {
			h.logger.Info("Token expired", zap.String("token", req.Token))
			return nil, ErrTokenExpired
		}
		h.logger.Error("Error sessionsDB.GetUserIDByToken", zap.Error(err))
		return nil, ErrInternalServer
	}

	var userInfo *db.User
	if userInfo, err = h.usersDB.SelectByUserID(ctx, userID); err != nil {
		h.logger.Error("Error usersDB.SelectByUserID", zap.Error(err))
		if err == db.NotFound {
			return nil, ErrAccountNotFound
		}
		return nil, ErrInternalServer
	}

	var tokenUserID int64
	if tokenUserID, err = h.getUserIDFromToken(userInfo, []byte(req.Token)); err != nil {
		return nil, err
	}

	if tokenUserID != userID {
		return nil, ErrInvalidToken
	}

	return &proto.CheckTokenResponse{Token: req.Token}, nil
}

func (h *handler) getUserIDFromToken(userInfo *db.User, payload []byte) (_ int64, err error) {
	var block *pem.Block
	block, _ = pem.Decode([]byte(userInfo.Password))
	if block == nil || block.Type != rsaPrivateKey {
		h.logger.Error("Error pem.Decode: wrong block type")
		return 0, ErrInternalServer
	}

	var privateKey *rsa.PrivateKey
	if privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		h.logger.Error("Error x509.ParsePKCS1PrivateKey", zap.Error(err))
		return 0, ErrInternalServer
	}

	var token jwt.Token
	if token, err = jwt.Parse(
		payload,
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, &privateKey.PublicKey),
	); err != nil {
		h.logger.Error("Error jwt.Parse", zap.Error(err))
		return 0, ErrInvalidToken
	}

	var (
		value interface{}
		ok    bool
	)
	if value, ok = token.Get(db.UserIDField); !ok {
		h.logger.Error("Error token.Get UserIDField", zap.Error(err))
		return 0, ErrInvalidToken
	}

	var cast float64
	if cast, ok = value.(float64); !ok {
		h.logger.Error("Error cast UserIDField to float64", zap.Error(err))
		return 0, ErrInvalidToken
	}

	return int64(cast), nil
}
