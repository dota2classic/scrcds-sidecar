package srcds

import (
	"log"
	"os"
	"time"

	"github.com/gorcon/rcon" // or your actual RCON client package
)

// RunHeartbeatPoller periodically checks if the game server responds.
// If the server fails more than maxFails times in a row, uploadAndExit() is called.
func RunHeartbeatPoller() {
	const (
		addr     = "localhost:27015"
		password = "rconpassword"
		interval = 1 * time.Second
		maxFails = 15
	)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	consecutiveFails := 0

	for range ticker.C {
		if pollServer(addr, password) {
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
func pollServer(addr, password string) bool {
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
	UploadArtifacts()
	os.Exit(0)
}
