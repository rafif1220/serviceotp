package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"), // Kosongin aja kalau nggak diset di compose
		DB:       0,  // Pakai database default (0)
	})

	// Tes Ping ke Redis
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Gagal connect ke Redis: %v", err)
	}

	log.Println("Redis Connected!")
	return client
}