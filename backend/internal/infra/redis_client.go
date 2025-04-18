package infra

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redisのアドレス
		Password: "",               // パスワード（なければ空文字）
		DB:       0,                // デフォルトDB
	})
}
