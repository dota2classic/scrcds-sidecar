package main

import (
	"log"
	"net"
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
	net.DefaultResolver.PreferGo = true
	log.Println("IPv4 Forced")

	state.InitGlobalState()
	s3.InitS3Client()

	redis.InitRedisClient()
	rabbit.InitRabbitPublisher()

	if err := srcds.AwaitServerReady(1 * time.Minute); err != nil {
		log.Fatalf("Server never started: %v", err)
	}

	go srcds.RunHeartbeatPoller()

	redis.ServerStatus(true)

	// emit rmq
	rabbit.PublishSrcdsServerStartedEvent(&d2cmodels.SrcdsServerStartedEvent{
		MatchId: state.GlobalMatchInfo.MatchID,
		Server:  state.GlobalMatchInfo.ServerAddress,
	})

	http.Listen(7777)
}
