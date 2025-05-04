package main

import (
	"context"
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestCrawler(t *testing.T) {
	const url string = "https://aws.amazon.com/blogs/aws/aws-announces-pixtral-large-25-02-model-in-amazon-bedrock-serverless/"
	crawler(url)
}

func TestSaveToRedis(t *testing.T) {
	const url string = "https://aws.amazon.com/blogs/aws/aws-announces-pixtral-large-25-02-model-in-amazon-bedrock-serverless/"

	saveToRedis(url)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()

	val, err := client.Get(ctx, url).Result()
	if err != nil {
		panic(err)
	}

	log.Println(val)
}
