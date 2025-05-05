package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func saveToRedis(address string, content string, redis_url string) {
	client := redis.NewClient(&redis.Options{
		Addr:     redis_url,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()

	err := client.Set(ctx, address, content, 0).Err()
	if err != nil {
		log.Fatalf("Error saving to Redis: %s", err)
	}

}

func getFromRedis(address string, redis_url string) string {
	client := redis.NewClient(&redis.Options{
		Addr:     redis_url,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()

	val, err := client.Get(ctx, address).Result()
	if err != nil {
		log.Fatalf("Error getting from Redis: %s", err)
	}

	return val
}
