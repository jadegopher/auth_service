package sessions

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"auth/internal/infrastructure/db"
)

type sessionService struct {
	client *redis.Client
}

func New(client *redis.Client) *sessionService {
	return &sessionService{client: client}
}

func (s *sessionService) Insert(token string, userID int64, expiresAt time.Duration) (err error) {
	if _, err = s.client.Set(token, userID, expiresAt).Result(); err != nil {
		return err
	}

	return nil
}

func (s *sessionService) GetUserIDByToken(token string) (userID int64, err error) {
	var res string
	if res, err = s.client.Get(token).Result(); err != nil {
		if err == redis.Nil {
			return 0, db.NotFound
		}
		return 0, err
	}

	if userID, err = strconv.ParseInt(res, 10, 64); err != nil {
		return 0, err
	}

	return userID, nil
}
