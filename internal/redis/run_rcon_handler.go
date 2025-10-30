package redis

import (
	"log"
	"sidecar/internal/models"
	"sidecar/internal/srcds/rcon"
	"sidecar/internal/state"
)

func handleRunRcon(evt models.RunRconCommand) {
	log.Printf("Received RunRconCommand: %+v", evt)
	log.Printf("%s = %s", state.GlobalMatchInfo.ServerAddress, evt.ServerUrl)
	if state.GlobalMatchInfo.ServerAddress != evt.ServerUrl {
		return
	}

	// Try execute rcon
	response, err := rcon.RunRconCommand(evt.Command)
	if err != nil {
		log.Printf("Error running rcon command: %v", err)
	}
	log.Printf("Rcon command result %s", response)
}
