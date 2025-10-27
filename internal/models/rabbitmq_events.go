package models

type PlayerIdType struct {
	Value string `json:"value"`
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

type PlayerAbandonedEvent struct {
	PlayerId     PlayerIdType       `json:"playerId"`
	MatchId      int64              `json:"matchId"`
	AbandonIndex int                `json:"abandonIndex"`
	Mode         MatchmakingMode    `json:"mode"`
	GameState    DotaGameRulesState `json:"gameState"`
}
