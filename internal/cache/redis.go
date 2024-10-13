package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
	once        sync.Once
)

func InitializeRedisClient(redisAddr string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		panic(err)
	} else {
		fmt.Println("Connected to Redis successfully")
	}
}
