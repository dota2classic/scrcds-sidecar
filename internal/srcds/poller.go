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

// RunHeartbeatPoller periodically polls for metrics and checks that server is alive.
// Call this only after AwaitServerReady() has returned — the server is expected to
// be fully up, so any sustained failure means it has crashed.
// If the server fails more than maxFails times in a row, UploadAndExit() is called.
func RunHeartbeatPoller() {
	const (
		interval = 2 * time.Second
		maxFails = 25
	)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	consecutiveFails := 0

	for range ticker.C {
		if pollMetrics() {
			consecutiveFails = 0
		} else {
			consecutiveFails++
		}

		if consecutiveFails > maxFails {
			log.Println("Server became unresponsive — shutting down")
			UploadAndExit()
			return
		}
	}
}

// pollMetrics attempts to connect and run the RCON status command
func pollMetrics() bool {
	if err := metrics.CollectMetrics(); err != nil {
		return false
	}

	redis.ServerHeartbeat()

	return true
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
