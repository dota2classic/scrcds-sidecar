package models

type PlayerInfo struct {
	Hero        string  `json:"hero"`
	Level       int     `json:"level"`
	Bot         bool    `json:"bot"`
	PosX        float64 `json:"pos_x"`
	PosY        float64 `json:"pos_y"`
	Angle       float64 `json:"angle"`
	Mana        int     `json:"mana"`
	MaxMana     int     `json:"max_mana"`
	Health      int     `json:"health"`
	MaxHealth   int     `json:"max_health"`
	Item0       int     `json:"item0"`
	Item1       int     `json:"item1"`
	Item2       int     `json:"item2"`
	Item3       int     `json:"item3"`
	Item4       int     `json:"item4"`
	Item5       int     `json:"item5"`
	Kills       int     `json:"kills"`
	Deaths      int     `json:"deaths"`
	Assists     int     `json:"assists"`
	RespawnTime int     `json:"respawn_time"`
}

type SlotInfo struct {
	Team       int                 `json:"team"`
	SteamID    string              `json:"steam_id"`
	Connection DotaConnectionState `json:"connection"`
	HeroData   *PlayerInfo         `json:"hero_data,omitempty"`
}

type LiveMatchUpdateEvent struct {
	MatchID         int64              `json:"match_id"`
	MatchmakingMode MatchmakingMode    `json:"matchmaking_mode"`
	GameMode        DotaGameMode       `json:"game_mode"`
	GameState       DotaGameRulesState `json:"game_state"`
	Duration        int                `json:"duration"`
	Server          string             `json:"server"`
	Timestamp       int64              `json:"timestamp"`
	Towers          [2]int             `json:"towers"`
	Barracks        [2]int             `json:"barracks"`
	Heroes          []SlotInfo         `json:"heroes"`
}
