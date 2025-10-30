package main

import (
	"log"
	"sidecar/internal/http"
	"sidecar/internal/rabbit"
	"sidecar/internal/redis"
	"sidecar/internal/s3"
	"sidecar/internal/srcds"
	"sidecar/internal/state"
	"time"

	d2cmodels "github.com/dota2classic/d2c-go-models/models"
)

func main() {
	log.Printf("Initializing sidecar...")

	state.InitGlobalState()
	s3.InitS3Client()

	redis.InitRedisClient()
	rabbit.InitRabbitPublisher()

	go srcds.RunHeartbeatPoller()

	dur, _ := time.ParseDuration("1m")
	err := srcds.AwaitHeartbeat(dur)
	if err != nil {
		log.Fatalf("Server never started: %v", err)
	}

	redis.ServerStatus(true)

	// emit rmq
	rabbit.PublishSrcdsServerStartedEvent(&d2cmodels.SrcdsServerStartedEvent{
		MatchId: state.GlobalMatchInfo.MatchID,
		Server:  state.GlobalMatchInfo.ServerAddress,
	})

	http.Listen(7777)
}
