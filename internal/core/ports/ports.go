package ports

import (
	"context"
	"time"

	"auth/internal/core/entities"
)

type ISessions interface {
	Insert(token string, userID int64, expiresAt time.Duration) (err error)
	GetUserIDByToken(token string) (userID int64, err error)
}

type IUsers interface {
	Insert(ctx context.Context, name, password string) (userID int64, err error)
	SelectByUserID(ctx context.Context, userID int64) (user *entities.User, err error)
}
