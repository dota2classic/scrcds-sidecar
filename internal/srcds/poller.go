package srcds

import (
	"errors"
	"log"
	"os"
	"sidecar/internal/s3"
	"sidecar/internal/srcds/metrics"
	"sidecar/internal/state"
	"strconv"
	"time"

	"github.com/gorcon/rcon" // or your actual RCON client package
)

var hadSuccessfulHeartbeat bool = false

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
// If the server fails more than maxFails times in a row, uploadAndExit() is called.
func RunHeartbeatPoller() {
	rconPassword := os.Getenv("RCON_PASSWORD")
	const (
		interval = 1 * time.Second
		maxFails = 30
	)

	//addr := fmt.Sprintf("127.0.0.1:%d", state.GlobalMatchInfo.GameServerPort)
	addr := state.GlobalMatchInfo.ServerAddress
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	consecutiveFails := 0

	for range ticker.C {
		if pollMetrics(addr, rconPassword) {
			consecutiveFails = 0 // success
			hadSuccessfulHeartbeat = true
		} else {
			consecutiveFails++
		}

		if consecutiveFails > maxFails {
			log.Println("Server is unresponsive — shutting down")
			uploadAndExit()
			return
		}
	}
}

// pollMetrics attempts to connect and run the RCON status command
func pollMetrics(addr string, password string) bool {
	conn, err := rcon.Dial(addr, password)
	if err != nil {
		return false
	}
	defer conn.Close()

	metrics.CollectMetrics(conn)

	return true
}

// uploadAndExit is your shutdown handler — replace with your logic.
func uploadAndExit() {
	log.Println("Uploading files and exiting...")
	matchId, err := strconv.ParseInt(os.Getenv("MATCH_ID"), 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse MATCH_ID: %v", err)
	}
	s3.UploadArtifacts(matchId)

	os.Exit(0)
}
