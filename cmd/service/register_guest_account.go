package service

import (
	"auth/cmd/db"
	"auth/proto"
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

func (h *handler) RegisterGuestAccount(
	ctx context.Context,
	req *proto.RegisterGuestAccountRequest,
) (resp *proto.RegisterGuestAccountResponse, err error) {
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		h.logger.Error("Error rsa.GenerateKey", zap.Error(err))
		return nil, ErrInternalServer
	}

	var userID int64
	if userID, err = h.usersDB.Insert(
		ctx,
		req.Name,
		string(pem.EncodeToMemory(
			&pem.Block{
				Type:  rsaPrivateKey,
				Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
			},
		)),
	); err != nil {
		h.logger.Error("Error usersDB.Insert", zap.Error(err))
		return nil, ErrInternalServer
	}

	var token = jwt.New()
	if err = token.Set(jwt.IssuedAtKey, time.Now().UTC().Unix()); err != nil {
		h.logger.Error("Error set IssuedAtKey", zap.Error(err))
		return nil, ErrInternalServer
	}

	if err = token.Set(db.UserIDField, userID); err != nil {
		h.logger.Error("Error set UserIDField", zap.Error(err))
		return nil, ErrInternalServer
	}

	// Sign the token and generate a payload
	var signed []byte
	if signed, err = jwt.Sign(token, jwa.RS256, privateKey); err != nil {
		h.logger.Error("Error jwt.Sign", zap.Error(err))
		return nil, ErrInternalServer
	}

	if err = h.sessionsDB.Insert(string(signed), userID, time.Hour); err != nil {
		h.logger.Error("sessionsDB.Insert", zap.Error(err))
		return nil, ErrInternalServer
	}

	return &proto.RegisterGuestAccountResponse{Token: string(signed)}, nil
}
