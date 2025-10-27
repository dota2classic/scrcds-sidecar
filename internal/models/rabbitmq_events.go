package models

type PlayerIdType struct {
	value string
}

type PlayerConnectedEvent struct {
	PlayerId  PlayerIdType `json:"playerId"`
	MatchId   int64        `json:"matchId"`
	ServerUrl string       `json:"ServerUrl"`
	Ip        string       `json:"ip"`
}

type MatchFailedEvent struct {
	MatchID       int64
	Server        string
	FailedPlayers []string
	GoodParties   []string
}
