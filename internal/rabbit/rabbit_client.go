package rabbit

import (
	"log"
	"sidecar/internal/models"
)

func PublishMatchFailedEvent(event *models.MatchFailedEvent) {
	err := publishWithRetry(event, "MatchFailedEvent", 3)
	if err != nil {
		log.Println("Error publishing event event:", err)
		return
	}
	log.Println("Published MatchFailedEvent to rmq")
}

func PublishPlayerAbandonEvent(event *models.PlayerAbandonedEvent) {
	err := publishWithRetry(event, "PlayerAbandonEvent", 3)
	if err != nil {
		log.Println("Error publishing event event:", err)
		return
	}
	log.Println("Published PlayerAbandonedEvent to rmq")
}

func PublishGameResultsEvent(event *models.GameResultsEvent) {
	err := publishWithRetry(event, "GameResultsEvent", 3)
	if err != nil {
		log.Println("Error publishing event event:", err)
		return
	}
	log.Println("Published GameResultsEvent to rmq")
}
