package main

import (
	"log"
	"sidecar/internal/http"
	"sidecar/internal/rabbit"
	"sidecar/internal/redis"
	"sidecar/internal/s3"
	"sidecar/internal/srcds"
	"sidecar/internal/state"
)

func main() {
	log.Printf("Initializing sidecar...")

	state.InitGlobalState()
	s3.InitS3Client()

	redis.InitRedisClient()
	rabbit.InitRabbitPublisher()

	go srcds.RunHeartbeatPoller()

	http.Listen(7777)
}
