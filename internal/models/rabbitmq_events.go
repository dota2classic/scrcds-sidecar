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

type PlayerInMatchDTO struct {
	SteamID             string `json:"steam_id"`
	Team                int    `json:"team"`
	Kills               int    `json:"kills"`
	Deaths              int    `json:"deaths"`
	Assists             int    `json:"assists"`
	Level               int    `json:"level"`
	Item0               int    `json:"item0"`
	Item1               int    `json:"item1"`
	Item2               int    `json:"item2"`
	Item3               int    `json:"item3"`
	Item4               int    `json:"item4"`
	Item5               int    `json:"item5"`
	GPM                 int    `json:"gpm"`
	XPM                 int    `json:"xpm"`
	LastHits            int    `json:"last_hits"`
	Denies              int    `json:"denies"`
	Networth            int    `json:"networth"`
	HeroDamage          int    `json:"heroDamage"`
	TowerDamage         int    `json:"towerDamage"`
	HeroHealing         int    `json:"heroHealing"`
	Abandoned           bool   `json:"abandoned"`
	Hero                string `json:"hero"`
	PartyIndex          int    `json:"partyIndex"`
	SupportAbilityValue int    `json:"supportAbilityValue"`
	SupportGold         int    `json:"supportGold"`
	Misses              int    `json:"misses"`
	Bear                []int  `json:"bear,omitempty"`
}

type GameResultsEvent struct {
	MatchID         int64              `json:"matchId"`
	Winner          DotaTeam           `json:"winner"`
	Duration        int                `json:"duration"`
	GameMode        DotaGameMode       `json:"gameMode"`
	Type            MatchmakingMode    `json:"type"`
	Timestamp       int64              `json:"timestamp"`
	Server          string             `json:"server"`
	Patch           DotaPatch          `json:"patch"`
	Region          Region             `json:"region"`
	Players         []PlayerInMatchDTO `json:"players"`
	TowerStatus     []int              `json:"towerStatus"`
	BarracksStatus  []int              `json:"barracksStatus"`
	ExternalMatchID *int64             `json:"externalMatchId,omitempty"`
}
