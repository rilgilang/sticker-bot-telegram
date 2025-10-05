package db

import (
	"github.com/redis/go-redis/v9"
	"tele-sticker-finder/config"
)

// simple db connection
func NewRedisConnection(config *config.Config) *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "guest", // no password set
		DB:       0,       // use default DB
	})
}
