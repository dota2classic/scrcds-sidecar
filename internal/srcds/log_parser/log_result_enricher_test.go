package log_parser

import (
	"os"
	"sidecar/internal/models"
	"testing"
	"time"
)

func TestFillAdditionalDataFromLog_1x1(t *testing.T) {
	// Given
	evt := models.GameResultsEvent{
		MatchID:        42,
		Winner:         models.DOTA_TEAM_RADIANT,
		Duration:       420,
		GameMode:       models.DOTA_GAME_MODE_ALLPICK,
		Type:           models.MATCHMAKING_MODE_SOLOMID,
		Timestamp:      time.Now().UnixMilli(),
		Server:         "fsdfdsf:4242",
		Region:         models.REGION_RU_MOSCOW,
		Patch:          models.PATCH_DOTA_684,
		BarracksStatus: []int{},
		TowerStatus:    []int{},
		Players: []models.PlayerInMatchDTO{
			{
				SteamID:             "1608039572",
				Team:                int(models.DOTA_TEAM_RADIANT),
				Kills:               1,
				Deaths:              2,
				Assists:             3,
				Level:               4,
				Item0:               1,
				SupportAbilityValue: 900,
				SupportGold:         75,
				Misses:              10,
				Hero:                "npc_dota_hero_nevermore",
			},
			{
				SteamID:             "1139947395",
				Team:                int(models.DOTA_TEAM_RADIANT),
				Kills:               1,
				Deaths:              2,
				Assists:             3,
				Level:               4,
				Item0:               1,
				SupportAbilityValue: 186,
				SupportGold:         0,
				Misses:              44,
				Hero:                "npc_dota_hero_tinker",
			},
		},
	}

	// When
	logFile := "testdata/1x1.log"
	if _, err := os.Stat(logFile); err != nil {
		t.Fatalf("missing test log file: %v", err)
	}

	if err := FillAdditionalDataFromLog(&evt, logFile); err != nil {
		t.Fatalf("failed to fill additional data: %v", err)
	}

	// Then
	p0 := evt.Players[0]
	p1 := evt.Players[1]

	if p0.TowerDamage != 62 {
		t.Errorf("expected player[0].towerDamage = 62, got %v", p0.TowerDamage)
	}
	if p0.GPM != 310 {
		t.Errorf("expected player[0].gpm = 310, got %v", p0.GPM)
	}
	if p0.Networth != 4435 {
		t.Errorf("expected player[0].networth = 4435, got %v", p0.Networth)
	}
	if p0.Misses != 10 {
		t.Errorf("expected player[0].misses = 10, got %v", p0.Misses)
	}
	if p1.Misses != 44 {
		t.Errorf("expected player[1].misses = 44, got %v", p1.Misses)
	}

	expectedBear := []int{50, 182, 172, 143, 0, 0}
	if len(p0.Bear) != len(expectedBear) {
		t.Errorf("expected player[0].bear = %v, got %v", expectedBear, p0.Bear)
	} else {
		for i, v := range expectedBear {
			if p0.Bear[i] != v {
				t.Errorf("expected player[0].bear[%d] = %v, got %v", i, v, p0.Bear[i])
			}
		}
	}

	if p1.Bear != nil {
		t.Errorf("expected player[1].bear = nil, got %v", p1.Bear)
	}

	expectedTowerStatus := []int{2047, 2039}
	expectedBarracksStatus := []int{63, 63}

	if !equalIntSlice(evt.TowerStatus, expectedTowerStatus) {
		t.Errorf("expected towerStatus = %v, got %v", expectedTowerStatus, evt.TowerStatus)
	}
	if !equalIntSlice(evt.BarracksStatus, expectedBarracksStatus) {
		t.Errorf("expected barracksStatus = %v, got %v", expectedBarracksStatus, evt.BarracksStatus)
	}
}

func equalIntSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
