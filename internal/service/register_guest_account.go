package service

import (
	"context"
	"go.uber.org/zap"
	"time"
)

const (
	tokenLifeTime = time.Hour
	issuedAtField = "iat"
	userIDField   = "user_id"
)

// RegisterGuestAccount needs for register temp guest accounts
func (s *service) RegisterGuestAccount(ctx context.Context, username string) (string, error) {
	// Generate key
	key, err := s.keyGenerator.Generate()
	if err != nil {
		return "", err
	}

	// Store username and private key to database
	// TODO add transaction and if some error rollback it
	userID, err := s.usersDB.Insert(ctx, username, key.String())
	if err != nil {
		s.logger.Error("error usersDB.Insert", zap.Error(err))
		return "", err
	}

	// Create token for user
	token, err := s.tokenCreator.Create(
		map[string]interface{}{
			issuedAtField: time.Now().UTC().Unix(),
			userIDField:   userID,
		},
		key.Interface(),
	)
	if err != nil {
		return "", err
	}

	// Store token in database
	if err = s.sessionsDB.Insert(token, userID, tokenLifeTime); err != nil {
		s.logger.Error("error sessionsDB.Insert", zap.Error(err))
		return "", err
	}

	return token, nil
}
