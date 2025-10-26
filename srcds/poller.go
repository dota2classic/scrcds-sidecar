package srcds

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gorcon/rcon" // or your actual RCON client package
)

// RunHeartbeatPoller periodically checks if the game server responds.
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
		if pollServer(addr, rconPassword) {
			consecutiveFails = 0 // success
		} else {
			consecutiveFails++
			log.Printf("Heartbeat failed (%d/%d)", consecutiveFails, maxFails)
		}

		if consecutiveFails > maxFails {
			log.Println("Server unresponsive — triggering uploadAndExit()")
			uploadAndExit()
			return
		}
	}
}

// pollServer attempts to connect and run the RCON status command
func pollServer(addr string, password string) bool {
	conn, err := rcon.Dial(addr, password)
	if err != nil {
		log.Printf("Failed to connect to RCON: %v", err)
		return false
	}
	defer conn.Close()

	_, err = conn.Execute("status")
	if err != nil {
		log.Printf("Failed to execute RCON command: %v", err)
		return false
	}

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
