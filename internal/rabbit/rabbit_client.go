package rabbit

import (
	"log"
	"sidecar/internal/models"
)

func PublishMatchFailedEvent(event models.MatchFailedEvent) {
	err := publishWithRetry(event, 3)
	log.Println("Published MatchFailedEvent to rmq")
	if err != nil {
		log.Println("Error publishing event event:", err)
		return
	}
}
