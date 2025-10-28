package models

import (
	"strconv"

	d2cmodels "github.com/dota2classic/d2c-go-models/models"
)

func ParseLobbyType(raw string) d2cmodels.MatchmakingMode {
	val, _ := strconv.Atoi(raw)
	return d2cmodels.MatchmakingMode(val)
}

func ParseGameMode(raw string) d2cmodels.DotaGameMode {
	val, _ := strconv.Atoi(raw)
	return d2cmodels.DotaGameMode(val)
}
