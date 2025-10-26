package srcds

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gorcon/rcon" // or your actual RCON client package
)

// RunHeartbeatPoller periodically polls for metrics and checks that server is alive
// If the server fails more than maxFails times in a row, uploadAndExit() is called.
func RunHeartbeatPoller() {
	rconPassword := os.Getenv("RCON_PASSWORD")
	const (
		addr     = "127.0.0.1:27015"
		interval = 1 * time.Second
		maxFails = 30
	)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	consecutiveFails := 0

	for range ticker.C {
		if pollMetrics(addr, rconPassword) {
			consecutiveFails = 0 // success
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
		log.Printf("Failed to connect to RCON: %v", err)
		return false
	}
	defer conn.Close()

	stats, err := conn.Execute("stats")
	if err != nil {
		log.Printf("Failed to execute RCON command: %v", err)
		return false
	}

	_, err = ParseStatsResponse(stats)

	return true
}

// uploadAndExit is your shutdown handler — replace with your logic.
func uploadAndExit() {
	log.Println("Uploading files and exiting...")
	matchId, err := strconv.ParseInt(os.Getenv("MATCH_ID"), 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse MATCH_ID: %v", err)
	}
	UploadArtifacts(matchId)

	os.Exit(0)
}
