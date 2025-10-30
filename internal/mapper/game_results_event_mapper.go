package mapper

import (
	"sidecar/internal/models"
	"sidecar/internal/state"
	"sidecar/internal/util/dotamaps"
	"strconv"

	d2cmodels "github.com/dota2classic/d2c-go-models/models"
)

func MapGameResults(d models.MatchFinishedOnSRCDS) d2cmodels.GameResultsEvent {
	// Collect unique party IDs
	parties := uniqueParties(d.Players)

	// Map players
	players := make([]d2cmodels.PlayerInMatchDTO, 0, len(d.Players))
	for _, p := range d.Players {
		partyIndex := indexOfParty(parties, p.PartyID) + 1 // +1 same as TS version

		player := d2cmodels.PlayerInMatchDTO{
			SteamID:    strconv.FormatInt(p.SteamID, 10),
			PartyIndex: partyIndex,
			Team:       p.Team,
			Kills:      p.Kills,
			Deaths:     p.Deaths,
			Assists:    p.Assists,
			Level:      p.Level,

			Item0: itemIDByName(p.Items, 0),
			Item1: itemIDByName(p.Items, 1),
			Item2: itemIDByName(p.Items, 2),
			Item3: itemIDByName(p.Items, 3),
			Item4: itemIDByName(p.Items, 4),
			Item5: itemIDByName(p.Items, 5),

			GPM:       int(p.GPM),
			XPM:       int(p.XPM),
			LastHits:  int(p.LastHits),
			Denies:    int(p.Denies),
			Abandoned: p.Connection == d2cmodels.DOTA_CONNECTION_STATE_ABANDONED || p.Connection == d2cmodels.DOTA_CONNECTION_STATE_FAILED,
			Networth:  int(p.Networth),

			HeroDamage:  0,
			HeroHealing: 0,
			TowerDamage: 0,

			Hero: p.Hero,

			SupportAbilityValue: 0,
			SupportGold:         0,
			Misses:              0,
		}

		players = append(players, player)
	}

	return d2cmodels.GameResultsEvent{
		MatchID:        d.MatchID,
		Winner:         d.Winner,
		Duration:       d.Duration,
		GameMode:       d.GameMode,
		Type:           d.Type,
		Timestamp:      d.Timestamp,
		Server:         state.GlobalMatchInfo.ServerAddress,
		Patch:          d.Patch,
		Region:         d.Region,
		Players:        players,
		TowerStatus:    []int{0, 0},
		BarracksStatus: []int{0, 0},
	}
}

// --- Helpers ---

func uniqueParties(players []models.SRCDSPlayer) []string {
	m := map[string]struct{}{}
	var result []string
	for _, p := range players {
		if _, exists := m[p.PartyID]; !exists {
			m[p.PartyID] = struct{}{}
			result = append(result, p.PartyID)
		}
	}
	return result
}

func indexOfParty(parties []string, id string) int {
	for i, p := range parties {
		if p == id {
			return i
		}
	}
	return -1
}

func itemIDByName(items []string, idx int) int {
	if idx >= len(items) || items[idx] == "" {
		return 0
	}
	return dotamaps.ItemID(items[idx]) // assumes you have a helper that maps names to IDs
}
