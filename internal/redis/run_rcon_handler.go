package redis

import (
	"log"
	"sidecar/internal/models"
	"sidecar/internal/srcds/rcon"
	"sidecar/internal/state"
)

func handleRunRcon(evt models.RunRconCommand) {
	if state.GlobalMatchInfo.ServerAddress != evt.ServerUrl {
		return
	}

	// Try execute rcon
	_, err := rcon.RunRconCommand(evt.Command)
	if err != nil {
		log.Printf("Error running rcon command: %v", err)
	}
}
