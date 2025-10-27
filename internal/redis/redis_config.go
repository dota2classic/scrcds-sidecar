package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var client *redis.Client

func InitRedisClient() {
	// Create client
	client = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379", // or your Redis service in k8s
		Password:     "",               // set if required
		DB:           0,
		PoolSize:     10, // adjust for concurrency
		MinIdleConns: 2,  // keep a few warm
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Redis client initialized")
}
