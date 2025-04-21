package infra

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient() *redis.Client {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("invalid redis URL: %v", err)
	}
	return redis.NewClient(opt)

}
