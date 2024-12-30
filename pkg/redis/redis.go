package redis

import (
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedis(options *redis.Options) *redis.Client {
	log.Printf("Starting Redis Connection")
	return redis.NewClient(options)
}
