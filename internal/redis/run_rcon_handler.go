package redis

import (
	"log"
	"sidecar/internal/models"
	"sidecar/internal/srcds/rcon"
	"sidecar/internal/state"
)

func handleRunRcon(evt *models.RunRconCommand) (*models.RunRconResponse, error) {
	log.Printf("Received RunRconCommand: %+v", evt)
	if state.GlobalMatchInfo.ServerAddress != evt.ServerUrl {
		return nil, nil
	}

	// Try execute rcon
	response, err := rcon.RunRconCommand(evt.Command)
	if err != nil {
		log.Printf("Error running rcon command: %v", err)
	}
	log.Printf("Rcon command result %s", response)
	return &models.RunRconResponse{Response: response}, nil
}
