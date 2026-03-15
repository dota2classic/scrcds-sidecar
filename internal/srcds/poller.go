package srcds

import (
	"log"
	"os"
	"sidecar/internal/redis"
	"sidecar/internal/s3"
	"sidecar/internal/srcds/metrics"
	"strconv"
	"time"
)

// RunHeartbeatPoller periodically checks server health via A2S and collects metrics.
// Call this only after AwaitServerReady() has returned.
// Health/upness is determined by A2S responses — RCON metric failures are logged
// but do not count as the server being down.
// If A2S fails more than maxFails times in a row, UploadAndExit() is called.
func RunHeartbeatPoller() {
	const (
		interval = 2 * time.Second
		maxFails = 25
	)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	consecutiveFails := 0

	for range ticker.C {
		if isServerReady() {
			consecutiveFails = 0
			if err := metrics.CollectMetrics(); err != nil {
				log.Printf("Metrics collection failed (server still up): %v", err)
			}
			redis.ServerHeartbeat()
		} else {
			consecutiveFails++
			log.Printf("A2S health check failed (%d/%d)", consecutiveFails, maxFails)
		}

		if consecutiveFails > maxFails {
			log.Println("Server became unresponsive — shutting down")
			UploadAndExit()
			return
		}
	}
}

// UploadAndExit is your shutdown handler — replace with your logic.
func UploadAndExit() {
	log.Println("Uploading files and exiting...")
	matchId, err := strconv.ParseInt(os.Getenv("MATCH_ID"), 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse MATCH_ID: %v", err)
	}
	s3.UploadArtifacts(matchId)

	redis.ServerStatus(false)

	metrics.Delete()
	os.Exit(0)
}
