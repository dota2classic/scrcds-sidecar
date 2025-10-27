package models

import "strconv"

func ParseLobbyType(raw string) MatchmakingMode {
	val, _ := strconv.Atoi(raw)
	return MatchmakingMode(val)
}

func ParseGameMode(raw string) DotaGameMode {
	val, _ := strconv.Atoi(raw)
	return DotaGameMode(val)
}
