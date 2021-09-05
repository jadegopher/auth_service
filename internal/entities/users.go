package entities

import (
	"context"
	"time"
)

type IUsers interface {
	Insert(ctx context.Context, name, password string) (userID int64, err error)
	SelectByUserID(ctx context.Context, userID int64) (user *User, err error)
}

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
