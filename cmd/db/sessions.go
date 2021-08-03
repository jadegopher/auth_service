package db

import "time"

const UserIDField = "user_id"

type ISessions interface {
	Insert(token string, userID int64, expiresAt time.Duration) (err error)
	GetUserIDByToken(token string) (userID int64, err error)
}
