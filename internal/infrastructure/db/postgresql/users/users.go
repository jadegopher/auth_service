package users

import (
	"auth/internal/entities"
	"auth/internal/infrastructure/db"
	"context"
	"database/sql"
)

type userService struct {
	conn *sql.DB
}

func New(conn *sql.DB) *userService {
	return &userService{conn: conn}
}

func (u *userService) Insert(ctx context.Context, name string, password string) (userID int64, err error) {
	queryInsert := `
		INSERT INTO users (name, password)
		VALUES ($1, $2)
		RETURNING id;`
	if err = u.conn.QueryRowContext(ctx, queryInsert, name, password).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

func (u *userService) SelectByUserID(ctx context.Context, userID int64) (_ *entities.User, err error) {
	querySelect := `
		SELECT id, name, password, created_at, updated_at
		FROM users
		WHERE id = $1;
	`
	user := &entities.User{}
	if err = u.conn.QueryRowContext(ctx, querySelect, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.NotFound
		}
		return nil, err
	}
	return user, nil
}
