package dotamaps

import "testing"

func TestHeroID(t *testing.T) {
	heroId := HeroID("npc_dota_hero_terrorblade")
	if heroId != 109 {
		t.Errorf("HeroID was wrong, expected %v, got %v", 109, heroId)
	}
}
