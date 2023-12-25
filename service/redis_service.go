package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type StorageService struct {
	redisClient *redis.Client
}

var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

const CacheDuration = 6 * time.Hour

func InitializeStore(addr string, password string, db int) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to init redis: %v", err)
	}

	fmt.Printf("Redis started successfully, ping message = %s\n", pong)
	storeService.redisClient = redisClient
}

func CloseStoreRedisConn() {
	storeService.redisClient.Close()
}

func SaveUrlMapping(shortUrl string, originUrl string, userId string) error {
	return storeService.redisClient.Set(ctx, shortUrl, originUrl, CacheDuration).Err()
}

func RetrieveInitialUrl(shortUrl string) (string, error) {
	result, err := storeService.redisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}
