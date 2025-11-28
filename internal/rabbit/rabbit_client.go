package rabbit

import (
	"log"
	"sidecar/internal/models"

	d2cmodels "github.com/dota2classic/d2c-go-models/models"
)

func PublishMatchFailedEvent(event *d2cmodels.MatchFailedEvent) {
	err := publishWithRetry(event, "MatchFailedEvent", 3)
	if err != nil {
		log.Println("Error publishing event event:", err)
		return
	}
	log.Println("Published MatchFailedEvent to rmq")
}

func PublishPlayerAbandonEvent(event *d2cmodels.PlayerAbandonedEvent) {
	err := publishWithRetry(event, "PlayerAbandonedEvent", 3)
	if err != nil {
		log.Println("Error publishing event event:", err)
		return
	}
	log.Println("Published PlayerAbandonedEvent to rmq")
}

func PublishGameResultsEvent(event *d2cmodels.GameResultsEvent) {
	err := publishWithRetry(event, "GameResultsEvent", 3)
	if err != nil {
		log.Println("Error publishing event event:", err)
		return
	}
	log.Println("Published GameResultsEvent to rmq")
}

func PublishSrcdsServerStartedEvent(event *d2cmodels.SrcdsServerStartedEvent) {
	err := publishWithRetry(event, "SrcdsServerStartedEvent", 3)
	if err != nil {
		log.Println("Error publishing event event:", err)
		return
	}
	log.Println("Published SrcdsServerStartedEvent to rmq")
}

func PublishArtifactUploadedEvent(event *models.MatchArtifactUploadedEvent) {
	err := publishWithRetry(event, "MatchArtifactUploadedEvent", 3)
	if err != nil {
		log.Println("Error publishing event event:", err)
		return
	}
	log.Println("Published MatchArtifactUploadedEvent to rmq")
}
