package service

import (
	"context"

	"go.uber.org/zap"

	"auth/internal/infrastructure/db"
)

// CheckToken needs for check is token valid or expired
func (s *Service) CheckToken(_ context.Context, token string) (ok bool, err error) {
	// Check token in database; if exists - everything is ok
	if _, err = s.sessionsDB.GetUserIDByToken(token); err != nil {
		if err == db.NotFound {
			return false, nil
		}
		s.logger.Error("error sessionsDB.GetUserIDByToken", zap.Error(err))
		return false, err
	}

	// TODO auto update existing token

	return true, nil
}
