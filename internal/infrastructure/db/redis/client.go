package redis

import (
	"fmt"

	"github.com/go-redis/redis"

	"auth/internal/core/entities"
)

func NewConnection(database entities.Database) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", database.IP, database.Port),
		Password: database.Password,
		DB:       0,
	})

	if _, err = client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}
