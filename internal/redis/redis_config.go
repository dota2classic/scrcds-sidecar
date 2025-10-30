package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	d2cutils "github.com/dota2classic/d2c-go-models/util"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var client *redis.Client

func InitRedisClient() {
	host := os.Getenv("REDIS_HOST")
	port := d2cutils.GetEnvInt("REDIS_PORT", 6379)

	password := os.Getenv("REDIS_PASSWORD")

	// Create client
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Password:     password,
		DB:           0,
		PoolSize:     2,
		MinIdleConns: 1,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Redis client initialized")

	Subscribe(context.Background(), client, "RunRconCommand", handleRunRcon)
}
