package rabbit

import (
	"log"
	"sidecar/internal/models"
)

func PublishMatchFailedEvent(event models.MatchFailedEvent) {
	err := publishWithRetry(event, 3)
	if err != nil {
		log.Println("Error publishing event event:", err)
		return
	}
}
