package models

import "github.com/dota2classic/d2c-go-models/models"

type ArtifactType string

const (
	ARTIFACT_TYPE_REPLAY ArtifactType = "REPLAY"
	ARTIFACT_TYPE_LOG    ArtifactType = "LOG"
)

type RunRconCommand struct {
	Command   string `json:"command"`
	ServerUrl string `json:"serverUrl"`
}

type RunRconResponse struct {
	Response string `json:"response"`
}

type ArtifactUploadedEvent struct {
	MatchID      int64                  `json:"matchId"`
	LobbyType    models.MatchmakingMode `json:"lobbyType"`
	ArtifactType ArtifactType           `json:"artifactType"`
	Bucket       string                 `json:"bucket"`
	Key          string                 `json:"key"`
}
