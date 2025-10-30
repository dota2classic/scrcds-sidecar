package models

import "github.com/dota2classic/d2c-go-models/models"

type HeroData struct {
	Bot         bool     `json:"bot"`
	PosX        float64  `json:"pos_x"`
	PosY        float64  `json:"pos_y"`
	Angle       float64  `json:"angle"`
	Hero        string   `json:"hero"`
	Level       int      `json:"level"`
	Health      float64  `json:"health"`
	MaxHealth   float64  `json:"max_health"`
	Mana        float64  `json:"mana"`
	MaxMana     float64  `json:"max_mana"`
	RespawnTime float64  `json:"respawn_time"`
	RDuration   float64  `json:"r_duration"`
	Items       []string `json:"items"`
	Kills       float64  `json:"kills"`
	Deaths      float64  `json:"deaths"`
	Assists     float64  `json:"assists"`
}

type SlotInfoDto struct {
	Team       int                        `json:"team"`
	SteamID    int                        `json:"steam_id"`
	Connection models.DotaConnectionState `json:"connection"`
	HeroData   *HeroData                  `json:"hero_data,omitempty"`
}

type FailedPlayerInfo struct {
	SteamID    int64                      `json:"steam_id"`
	PartyID    string                     `json:"party_id,omitempty"`
	Connection models.DotaConnectionState `json:"connection"`
}

type SRCDSPlayer struct {
	Hero        string                     `json:"hero"`
	SteamID     int64                      `json:"steam_id"`
	Team        int                        `json:"team"`
	Level       int                        `json:"level"`
	Kills       int                        `json:"kills"`
	Deaths      int                        `json:"deaths"`
	Assists     int                        `json:"assists"`
	Connection  models.DotaConnectionState `json:"connection"`
	GPM         float64                    `json:"gpm"`
	XPM         float64                    `json:"xpm"`
	LastHits    float64                    `json:"last_hits"`
	Denies      float64                    `json:"denies"`
	TowerKills  float64                    `json:"tower_kills"`
	Networth    float64                    `json:"networth"`
	RoshanKills float64                    `json:"roshan_kills"`
	Items       []string                   `json:"items"`
	PartyID     string                     `json:"party_id"`
}

type LiveMatchDto struct {
	MatchID         int64                     `json:"match_id"`
	MatchmakingMode models.MatchmakingMode    `json:"matchmaking_mode"`
	GameMode        models.DotaGameMode       `json:"game_mode"`
	GameState       models.DotaGameRulesState `json:"game_state"`
	Timestamp       int64                     `json:"timestamp"`
	Duration        float64                   `json:"duration"`
	Server          string                    `json:"server"`
	Towers          [2]int                    `json:"towers"`
	Barracks        [2]int                    `json:"barracks"`
	Heroes          []SlotInfoDto             `json:"heroes"`
}

type MatchFailedOnSRCDS struct {
	Players []FailedPlayerInfo `json:"players"`
	MatchID int64              `json:"match_id"`
	Server  string             `json:"server"`
}

type PlayerAbandonOnSRCDS struct {
	MatchID      int64                     `json:"match_id"`
	SteamID      int64                     `json:"steam_id"`
	AbandonIndex int                       `json:"abandon_index"`
	Mode         models.MatchmakingMode    `json:"mode"`
	Server       string                    `json:"server"`
	GameState    models.DotaGameRulesState `json:"game_state"`
}

type PlayerNotLoadedOnSRCDS struct {
	MatchID int64                  `json:"match_id"`
	SteamID int64                  `json:"steam_id"`
	Mode    models.MatchmakingMode `json:"mode"`
	Server  string                 `json:"server"`
}

type PlayerConnectedOnSRCDS struct {
	MatchID      int64                     `json:"match_id"`
	SteamID      int64                     `json:"steam_id"`
	LobbyType    models.MatchmakingMode    `json:"lobby_type"`
	Server       string                    `json:"server"`
	IP           string                    `json:"ip"`
	GameState    models.DotaGameRulesState `json:"gameState"`
	FirstConnect bool                      `json:"firstConnect"`
	Duration     float64                   `json:"duration"`
}

type MatchFinishedOnSRCDS struct {
	MatchID   int64                  `json:"MatchId"`
	Winner    models.DotaTeam        `json:"winner"`
	Duration  int                    `json:"duration"`
	Type      models.MatchmakingMode `json:"type"`
	GameMode  models.DotaGameMode    `json:"gameMode"`
	Timestamp int64                  `json:"timestamp"`
	Server    string                 `json:"server"`
	Region    models.Region          `json:"region"`
	Patch     models.DotaPatch       `json:"patch"`
	Players   []SRCDSPlayer          `json:"players"`
}
