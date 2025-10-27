package models

type HeroData struct {
	Bot         bool     `json:"bot"`
	PosX        float64  `json:"pos_x"`
	PosY        float64  `json:"pos_y"`
	Angle       float64  `json:"angle"`
	Hero        string   `json:"hero"`
	Level       int      `json:"level"`
	Health      int      `json:"health"`
	MaxHealth   int      `json:"max_health"`
	Mana        int      `json:"mana"`
	MaxMana     int      `json:"max_mana"`
	RespawnTime int      `json:"respawn_time"`
	RDuration   int      `json:"r_duration"`
	Items       []string `json:"items"`
	Kills       int      `json:"kills"`
	Deaths      int      `json:"deaths"`
	Assists     int      `json:"assists"`
}

type SlotInfoDto struct {
	Team       int                 `json:"team"`
	SteamID    int                 `json:"steam_id"`
	Connection DotaConnectionState `json:"connection"`
	HeroData   *HeroData           `json:"hero_data,omitempty"`
}

type FailedPlayerInfo struct {
	SteamID    int64               `json:"steam_id"`
	PartyID    string              `json:"party_id,omitempty"`
	Connection DotaConnectionState `json:"connection"`
}

type SRCDSPlayer struct {
	Hero        string              `json:"hero"`
	SteamID     int64               `json:"steam_id"`
	Team        int                 `json:"team"`
	Level       int                 `json:"level"`
	Kills       int                 `json:"kills"`
	Deaths      int                 `json:"deaths"`
	Assists     int                 `json:"assists"`
	Connection  DotaConnectionState `json:"connection"`
	GPM         int                 `json:"gpm"`
	XPM         int                 `json:"xpm"`
	LastHits    int                 `json:"last_hits"`
	Denies      int                 `json:"denies"`
	TowerKills  int                 `json:"tower_kills"`
	Networth    int                 `json:"networth"`
	RoshanKills int                 `json:"roshan_kills"`
	Items       []string            `json:"items"`
	PartyID     string              `json:"party_id"`
}

type LiveMatchDto struct {
	MatchID         int64              `json:"match_id"`
	MatchmakingMode MatchmakingMode    `json:"matchmaking_mode"`
	GameMode        DotaGameMode       `json:"game_mode"`
	GameState       DotaGameRulesState `json:"game_state"`
	Timestamp       int64              `json:"timestamp"`
	Duration        int                `json:"duration"`
	Server          string             `json:"server"`
	Towers          [2]int             `json:"towers"`
	Barracks        [2]int             `json:"barracks"`
	Heroes          []SlotInfoDto      `json:"heroes"`
}

type MatchFailedOnSRCDS struct {
	Players []FailedPlayerInfo `json:"players"`
	MatchID int64              `json:"match_id"`
	Server  string             `json:"server"`
}

type PlayerAbandonOnSRCDS struct {
	MatchID      int64              `json:"match_id"`
	SteamID      int64              `json:"steam_id"`
	AbandonIndex int                `json:"abandon_index"`
	Mode         MatchmakingMode    `json:"mode"`
	Server       string             `json:"server"`
	GameState    DotaGameRulesState `json:"game_state"`
}

type PlayerNotLoadedOnSRCDS struct {
	MatchID int64           `json:"match_id"`
	SteamID int64           `json:"steam_id"`
	Mode    MatchmakingMode `json:"mode"`
	Server  string          `json:"server"`
}

type PlayerConnectedOnSRCDS struct {
	MatchID      int64              `json:"match_id"`
	SteamID      int64              `json:"steam_id"`
	LobbyType    MatchmakingMode    `json:"lobby_type"`
	Server       string             `json:"server"`
	IP           string             `json:"ip"`
	GameState    DotaGameRulesState `json:"gameState"`
	FirstConnect bool               `json:"firstConnect"`
	Duration     int                `json:"duration"`
}

type MatchFinishedOnSRCDS struct {
	MatchID   int64           `json:"MatchId"`
	Winner    DotaTeam        `json:"winner"`
	Duration  int             `json:"duration"`
	Type      MatchmakingMode `json:"type"`
	GameMode  DotaGameMode    `json:"gameMode"`
	Timestamp int64           `json:"timestamp"`
	Server    string          `json:"server"`
	Region    Region          `json:"region"`
	Patch     DotaPatch       `json:"patch"`
	Players   []SRCDSPlayer   `json:"players"`
}
