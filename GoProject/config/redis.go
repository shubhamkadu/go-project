package config

import (
	"github.com/go-redis/redis/v8"
	"context"
	"log"
)

var RedisClient *redis.Client

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis")
}

