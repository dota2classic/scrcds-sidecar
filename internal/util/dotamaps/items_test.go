package dotamaps

import "testing"

func TestItemID_valid(t *testing.T) {
	itemId := ItemID("lotus_orb")

	if itemId != 226 {
		t.Errorf("itemId should be 226")
	}
}
