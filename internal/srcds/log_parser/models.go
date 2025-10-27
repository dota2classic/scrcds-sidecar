package log_parser

type ParsedProtobufMessage struct {
	Duration                  float64                `json:"duration"`
	GoodGuysWin               bool                   `json:"good_guys_win"`
	Date                      int64                  `json:"date"`
	NumPlayers                []int                  `json:"num_players"`
	Teams                     []team                 `json:"teams"`
	TowerStatus               []int                  `json:"tower_status"`
	BarracksStatus            []int                  `json:"barracks_status"`
	Cluster                   int64                  `json:"cluster"`
	ServerAddr                string                 `json:"server_addr"`
	FirstBloodTime            float64                `json:"first_blood_time"`
	GameBalance               float64                `json:"game_balance"`
	AutomaticSurrender        bool                   `json:"automatic_surrender"`
	ServerVersion             int64                  `json:"server_version"`
	AverageNetworthDelta      float64                `json:"average_networth_delta"`
	NetworthDeltaMin10        float64                `json:"networth_delta_min10"`
	NetworthDeltaMin20        float64                `json:"networth_delta_min20"`
	MaximumLosingNetworthLead float64                `json:"maximum_losing_networth_lead"`
	AverageExperienceDelta    float64                `json:"average_experience_delta"`
	ExperienceDeltaMin10      float64                `json:"experience_delta_min10"`
	ExperienceDeltaMin20      float64                `json:"experience_delta_min20"`
	BonusGoldWinnerMin10      float64                `json:"bonus_gold_winner_min10"`
	BonusGoldWinnerMin20      float64                `json:"bonus_gold_winner_min20"`
	BonusGoldWinnerTotal      float64                `json:"bonus_gold_winner_total"`
	BonusGoldLoserMin10       float64                `json:"bonus_gold_loser_min10"`
	BonusGoldLoserMin20       float64                `json:"bonus_gold_loser_min20"`
	BonusGoldLoserTotal       float64                `json:"bonus_gold_loser_total"`
	MatchID                   int64                  `json:"match_id"`
	RegionID                  int64                  `json:"region_id"`
	Players                   []playerConnectionInfo `json:"players"`
}

type team struct {
	Players []player `json:"players"`
}

type player struct {
	SteamID                  int64                     `json:"steam_id"`
	HeroID                   int                       `json:"hero_id"`
	Items                    []int64                   `json:"items"`
	Gold                     int64                     `json:"gold"`
	Kills                    int64                     `json:"kills"`
	Deaths                   int64                     `json:"deaths"`
	Assists                  int64                     `json:"assists"`
	LeaverStatus             int64                     `json:"leaver_status"`
	LastHits                 int64                     `json:"last_hits"`
	Denies                   int64                     `json:"denies"`
	GoldPerMin               int                       `json:"gold_per_min"`
	XpPerMinute              int                       `json:"xp_per_minute"`
	GoldSpent                int                       `json:"gold_spent"`
	Level                    int                       `json:"level"`
	HeroDamage               int                       `json:"hero_damage"`
	TowerDamage              int                       `json:"tower_damage"`
	HeroHealing              int                       `json:"hero_healing"`
	TimeLastSeen             int64                     `json:"time_last_seen"`
	SupportAbilityValue      int                       `json:"support_ability_value"`
	PartyID                  int64                     `json:"party_id"`
	ScaledKills              float64                   `json:"scaled_kills"`
	ScaledDeaths             float64                   `json:"scaled_deaths"`
	ScaledAssists            float64                   `json:"scaled_assists"`
	ClaimedFarmGold          int64                     `json:"claimed_farm_gold"`
	SupportGold              int                       `json:"support_gold"`
	ClaimedDenies            int                       `json:"claimed_denies"`
	ClaimedMisses            int                       `json:"claimed_misses"`
	Misses                   int                       `json:"misses"`
	AbilityUpgrades          []abilityUpgrade          `json:"ability_upgrades"`
	NetWorth                 int                       `json:"net_worth"`
	AdditionalUnitsInventory *additionalUnitsInventory `json:"additional_units_inventory,omitempty"`
}

type abilityUpgrade struct {
	Ability int64 `json:"ability"`
	Time    int64 `json:"time"`
}

type additionalUnitsInventory struct {
	UnitName string `json:"unit_name"`
	Items    []int  `json:"items"`
}

type playerConnectionInfo struct {
	AccountID     int64   `json:"account_id"`
	IP            int64   `json:"ip"`
	AvgPingMs     float64 `json:"avg_ping_ms"`
	PacketLoss    float64 `json:"packet_loss"`
	PingDeviation float64 `json:"ping_deviation"`
	FullResends   float64 `json:"full_resends"`
}
