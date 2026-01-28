package redis

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sidecar/internal/models"
	"time"

	d2cutils "github.com/dota2classic/d2c-go-models/util"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var client *redis.Client

func logDNS(host string) {
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Printf("DNS lookup failed for %s: %v", host, err)
		return
	}

	for _, ip := range ips {
		log.Printf("DNS resolved %s -> %s", host, ip.String())
	}
}

func InitRedisClient() {
	host := os.Getenv("REDIS_HOST")
	port := d2cutils.GetEnvInt("REDIS_PORT", 6379)
	password := os.Getenv("REDIS_PASSWORD")

	logDNS(host)

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Connecting to redis at %s", addr)

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
	for i := 1; i <= 10; i++ {
		if err := client.Ping(ctx).Err(); err != nil {
			log.Printf("redis ping failed (attempt %d): %v", i, err)
			time.Sleep(time.Second)
			continue
		}
		log.Println("Redis client initialized")
		break
	}

	Subscribe[models.RunRconCommand](context.Background(), client, "RunRconCommand", handleRunRcon)
}
