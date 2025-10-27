package log_parser

import (
	"fmt"
	"log"
	"os"
	"sidecar/internal/models"
	"sidecar/internal/util/dotamaps"
	"strconv"
)

// fillAdditionalDataFromLog enriches a GameResultsEvent using data from a log file.
func FillAdditionalDataFromLog(evt *models.GameResultsEvent, logFile string) error {
	log.Printf("Beginning parsing log file: %s for match %v", logFile, evt.MatchID)

	data, err := os.ReadFile(logFile)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	parsed, err := ParseLog(string(data))
	if err != nil {
		return fmt.Errorf("failed to parse log: %w", err)
	}

	evt.TowerStatus = parsed.TowerStatus
	evt.BarracksStatus = parsed.BarracksStatus

	for _, team := range parsed.Teams {
		for _, player := range team.Players {
			// Convert SteamID64 to SteamID32
			//steamID64 := new(big.Int)
			//steamID64.SetString(player.SteamID, 10)
			//
			//base := new(big.Int)
			//base.SetString("76561197960265728", 10)
			//steam32 := new(big.Int).Sub(steamID64, base).String()

			steam32 := 76561197960265728 - player.SteamID
			steam32str := strconv.FormatInt(steam32, 10)

			// Try to find player by Steam32
			var baseData *models.PlayerInMatchDTO
			for i := range evt.Players {
				if evt.Players[i].SteamID == steam32str {
					baseData = &evt.Players[i]
					break
				}
			}

			if baseData == nil {
				log.Printf("WARN: Didn't find base player data for steam id %d, trying hero id %v", player.SteamID, player.HeroID)
				// fallback: find by hero id
				for i := range evt.Players {
					if dotamaps.HeroID(evt.Players[i].Hero) == player.HeroID {
						baseData = &evt.Players[i]
						break
					}
				}
			}

			if baseData == nil {
				log.Printf("ERROR: Didn't find base player data for hero %v! Skipping player", player.HeroID)
				continue
			}

			// Enrich player data
			baseData.GPM = player.GoldPerMin
			baseData.XPM = player.XpPerMinute
			baseData.HeroDamage = player.HeroDamage
			baseData.HeroHealing = player.HeroHealing
			baseData.TowerDamage = player.TowerDamage
			baseData.Networth = player.NetWorth
			baseData.SupportGold = player.SupportGold
			baseData.SupportAbilityValue = player.SupportAbilityValue
			baseData.Misses = player.Misses
			if player.AdditionalUnitsInventory != nil {
				baseData.Bear = player.AdditionalUnitsInventory.Items
			}
		}
	}

	log.Printf("Log file parsed: %s for match %v", logFile, evt.MatchID)
	return nil
}
