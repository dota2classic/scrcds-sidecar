package mapper

import (
	"sidecar/internal/models"
	"sidecar/internal/util/dotamaps"
	"strconv"
)

func MapLiveMatchUpdatedEvent(dto models.LiveMatchDto) models.LiveMatchUpdateEvent {
	var mapped models.LiveMatchUpdateEvent
	mapped.MatchID = dto.MatchID
	mapped.Server = dto.Server
	mapped.MatchmakingMode = dto.MatchmakingMode
	mapped.GameMode = dto.GameMode
	mapped.GameState = dto.GameState
	mapped.Duration = dto.Duration
	mapped.Towers = dto.Towers
	mapped.Barracks = dto.Barracks
	mapped.Timestamp = dto.Timestamp

	for _, h := range dto.Heroes {
		var slot models.SlotInfo
		slot.Team = h.Team

		// SteamID might be number or string in JSON
		slot.SteamID = strconv.Itoa(h.SteamID)

		slot.Connection = h.Connection

		if h.HeroData != nil {
			hd := h.HeroData
			player := &models.PlayerInfo{
				Level:       hd.Level,
				Hero:        hd.Hero,
				Bot:         hd.Bot,
				PosX:        hd.PosX,
				PosY:        hd.PosY,
				Angle:       hd.Angle,
				Mana:        hd.Mana,
				MaxMana:     hd.MaxMana,
				Health:      hd.Health,
				MaxHealth:   hd.MaxHealth,
				Item0:       dotamaps.ItemID(hd.Items[0]),
				Item1:       dotamaps.ItemID(hd.Items[1]),
				Item2:       dotamaps.ItemID(hd.Items[2]),
				Item3:       dotamaps.ItemID(hd.Items[3]),
				Item4:       dotamaps.ItemID(hd.Items[4]),
				Item5:       dotamaps.ItemID(hd.Items[5]),
				Kills:       hd.Kills,
				Deaths:      hd.Deaths,
				Assists:     hd.Assists,
				RespawnTime: hd.RespawnTime,
			}
			slot.HeroData = player
		}

		mapped.Heroes = append(mapped.Heroes, slot)
	}

	return mapped
}
