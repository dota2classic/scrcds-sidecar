package log_parser

import (
	"os"
	"testing"
)

func TestParseLogFile_1x1(t *testing.T) {
	// Given
	data, _ := os.ReadFile("testdata/1x1.log")

	// When
	parsed, _ := ParseLog(string(data))

	// Then
	got1 := parsed.Teams[0].Players[0].TowerDamage
	got2 := parsed.Teams[1].Players[0].TowerDamage
	if got1 != 62 {
		t.Errorf("expected team[0].player[0].tower_damage = 62, got %v", got1)
	}
	if got2 != 0 {
		t.Errorf("expected team[1].player[0].tower_damage = 0, got %v", got2)
	}

	id1 := parsed.Teams[0].Players[0].SteamID
	id2 := parsed.Teams[1].Players[0].SteamID
	if id1 != "76561199568305300" {
		t.Errorf("expected team[0].player[0].steam_id = 76561199568305300, got %v", id1)
	}
	if id2 != "76561199100213123" {
		t.Errorf("expected team[1].player[0].steam_id = 76561199100213123, got %v", id2)
	}
}

func TestParseLogFile_5x5(t *testing.T) {
	// Given
	data, err := os.ReadFile("testdata/5x5.log")
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	// When
	parsed, err := ParseLog(string(data))
	if err != nil {
		t.Fatalf("failed to parse log: %v", err)
	}

	flattenPlayers := func(teams []team) []player {
		var all []player
		for _, t := range teams {
			all = append(all, t.Players...)
		}
		return all
	}

	players := flattenPlayers(parsed.Teams)

	// GPM
	expectedGPM := []int64{223, 432, 170, 374, 473, 359, 615, 278, 426, 592}
	for i, p := range players {
		if p.GoldPerMin != expectedGPM[i] {
			t.Errorf("expected player %d GPM = %v, got %v", i, expectedGPM[i], p.GoldPerMin)
		}
	}

	// XPM
	expectedXPM := []int64{229, 547, 179, 471, 549, 328, 610, 357, 487, 610}
	for i, p := range players {
		if p.XpPerMinute != expectedXPM[i] {
			t.Errorf("expected player %d XPM = %v, got %v", i, expectedXPM[i], p.XpPerMinute)
		}
	}

	// HeroDamage
	expectedHeroDamage := []int64{6243, 32146, 4324, 9731, 13355, 21939, 32220, 4320, 11411, 14658}
	for i, p := range players {
		if p.HeroDamage != expectedHeroDamage[i] {
			t.Errorf("expected player %d HeroDamage = %v, got %v", i, expectedHeroDamage[i], p.HeroDamage)
		}
	}

	// TowerDamage
	expectedTowerDamage := []int64{376, 0, 103, 242, 1079, 218, 4586, 355, 2240, 2023}
	for i, p := range players {
		if p.TowerDamage != expectedTowerDamage[i] {
			t.Errorf("expected player %d TowerDamage = %v, got %v", i, expectedTowerDamage[i], p.TowerDamage)
		}
	}

	// HeroHealing
	expectedHeroHealing := []int64{1371, 0, 0, 0, 0, 0, 0, 8377, 1624, 0}
	for i, p := range players {
		if p.HeroHealing != expectedHeroHealing[i] {
			t.Errorf("expected player %d HeroHealing = %v, got %v", i, expectedHeroHealing[i], p.HeroHealing)
		}
	}

	// Misses
	expectedMisses := []int64{14, 16, 16, 8, 10, 11, 22, 17, 18, 16}
	for i, p := range players {
		if p.Misses != expectedMisses[i] {
			t.Errorf("expected player %d Misses = %v, got %v", i, expectedMisses[i], p.Misses)
		}
	}

	// NetWorth
	expectedNetWorth := []int64{4562, 21169, 2409, 16635, 20675, 14937, 31282, 12210, 20559, 26524}
	for i, p := range players {
		if p.NetWorth != expectedNetWorth[i] {
			t.Errorf("expected player %d NetWorth = %v, got %v", i, expectedNetWorth[i], p.NetWorth)
		}
	}

	// Team sizes
	if len(parsed.Teams[0].Players) != 5 {
		t.Errorf("expected team[0] to have 5 players, got %d", len(parsed.Teams[0].Players))
	}
	if len(parsed.Teams[1].Players) != 5 {
		t.Errorf("expected team[1] to have 5 players, got %d", len(parsed.Teams[1].Players))
	}
}

func TestParseLogFile_IncompleteGame(t *testing.T) {
	// Given
	data, err := os.ReadFile("testdata/4x5.log")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	// When
	parsed, err := ParseLog(string(data))
	if err != nil {
		t.Fatalf("failed to parse log: %v", err)
	}

	// Tower status
	expectedTowerStatus := []int64{2047, 260}
	for i, v := range expectedTowerStatus {
		if parsed.TowerStatus[i] != v {
			t.Errorf("expected tower_status[%d] = %v, got %v", i, v, parsed.TowerStatus[i])
		}
	}

	// Barracks status
	expectedBarracksStatus := []int64{63, 51}
	for i, v := range expectedBarracksStatus {
		if parsed.BarracksStatus[i] != v {
			t.Errorf("expected barracks_status[%d] = %v, got %v", i, v, parsed.BarracksStatus[i])
		}
	}

	// Team players length
	if len(parsed.Teams[0].Players) != 5 {
		t.Errorf("expected team[0].players length = 5, got %d", len(parsed.Teams[0].Players))
	}
	if len(parsed.Teams[1].Players) != 5 {
		t.Errorf("expected team[1].players length = 5, got %d", len(parsed.Teams[1].Players))
	}
}

func TestParseLogFile_Druid(t *testing.T) {
	// Given
	data, err := os.ReadFile("testdata/druid.log")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	// When
	parsed, err := ParseLog(string(data))
	if err != nil {
		t.Fatalf("failed to parse log: %v", err)
	}

	inv := parsed.Teams[0].Players[1].AdditionalUnitsInventory

	if inv.UnitName != "spirit_bear" {
		t.Errorf("parsed.Teams[0].Players[1].AdditionalUnitsInventory.UnitName to be spirit bear, got %s", inv.UnitName)
	}

	expectedInventory := []int64{
		50,
		182,
		172,
		143,
		0,
		0,
	}

	for i := range expectedInventory {
		if expectedInventory[i] != inv.Items[i] {
			t.Errorf("Inventory item mismatch at index %d: expected %d got %d", i, expectedInventory[i], inv.Items[i])
		}
	}
}

func TestParseLogFile_Bots(t *testing.T) {
	// Given
	data, _ := os.ReadFile("testdata/bots.log")

	// When
	_, err := ParseLog(string(data))

	if err != nil {
		t.Errorf("ParseLog() returned nil")
	}
}
