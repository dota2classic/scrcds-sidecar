package srcds

import (
	"errors"
	"log"
	"os"
	"sidecar/internal/redis"
	"sidecar/internal/s3"
	"sidecar/internal/srcds/metrics"
	rcon2 "sidecar/internal/srcds/rcon"
	"strconv"
	"time"
)

var hadSuccessfulHeartbeat = false

// AwaitHeartbeat awaits for first successful heartbeat
func AwaitHeartbeat(maxWait time.Duration) error {
	ticker := time.NewTicker(1 * time.Second)
	start := time.Now()
	for range ticker.C {
		if time.Since(start) > maxWait {
			return errors.New("timed out waiting for heartbeat")
		}
		if hadSuccessfulHeartbeat {
			return nil
		}
		continue
	}
	defer ticker.Stop()
	return nil
}

// RunHeartbeatPoller periodically polls for metrics and checks that server is alive
// If the server fails more than maxFails times in a row, UploadAndExit() is called.
func RunHeartbeatPoller() {
	const (
		interval = 1 * time.Second
		maxFails = 5
	)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	consecutiveFails := 0

	for range ticker.C {
		if pollMetrics() {
			consecutiveFails = 0 // success
			hadSuccessfulHeartbeat = true
		} else {
			consecutiveFails++
		}

		if hadSuccessfulHeartbeat && consecutiveFails > maxFails {
			log.Println("Server became unresponsive — shutting down")
			UploadAndExit()
			return
		}
	}
}

// pollMetrics attempts to connect and run the RCON status command
func pollMetrics() bool {
	conn, err := rcon2.GetRconConnection()
	if err != nil {
		return false
	}
	defer conn.Close()

	metrics.CollectMetrics(conn)

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

	os.Exit(0)
}
