package srcds

import (
	"reflect"
	"testing"
)

func TestParseStatusResponse_BotGame(t *testing.T) {
	raw := `hostname: Dota 2
version : 41/41 0 secure  
steamid : [A:1:902557723:39451] (90241434204762139)
udp/ip  :  0.0.0.0:27015 os(Linux) type(dedicated)
sourcetv:  port 27020, delay 120.0s
players : 2 humans, 9 bots (15 max) (not hibernating)
edicts : 1037 used of 2048 max
gamestate: DOTA_GAMERULES_STATE_GAME_IN_PROGRESS Times: Transition=137.00 Current=1979.92
Lv Name         Player        K/ D/ A/ LH/ DN/ Gold Health    Mana
19 Unknown      фея те�  9/ 5/ 1/122/  7/  589    0/2259  294/ 702
12 Unknown      Негрос  1/10/ 4/ 89/  1/  610    0/ 986  436/ 572
12 Unknown      Edith Bot     2/ 8/ 6/ 15/  8/  644 1005/1005 1069/1069
11 Unknown      Vivian Bot    2/ 9/ 6/ 58/  4/  456 1138/1138  336/ 351
11 Unknown      Juan Bot      2/12/ 4/ 20/  1/  609 1518/1518  468/ 468
17 Unknown      Maurice Bot  11/ 3/12/117/ 40/  920  483/1613  184/ 780
16 Unknown      Chai Bot     11/ 2/ 8/ 56/  7/ 1297  752/1195  926/ 988
15 Unknown      Mordecai Bot  7/ 2/13/ 60/ 12/  420 1563/1613 1417/1417
13 Unknown      Jorge Bot     3/ 5/11/ 28/ 10/  180  929/ 929  806/ 806
16 Unknown      Lupe Bot      9/ 5/14/ 52/  3/ 1568 1095/1806  119/ 559
# userid name uniqueid connected ping loss state rate adr
# 2 "SourceTV" BOT active
# 3 "Edith Bot" BOT active
# 4 "Vivian Bot" BOT active
# 5 "Juan Bot" BOT active
# 6 "Maurice Bot" BOT active
# 7 "Chai Bot" BOT active
# 8 "Mordecai Bot" BOT active
# 9 "Jorge Bot" BOT active
#10 "Lupe Bot" BOT active
# 11 10 "Негросучка" [U:1:1037635443] 32:53 71 0 active 80000 123.31.25.81:12345
# 12 11 "фея технологий" [U:1:1175148404] 32:43 93 0 active 80000 123.219.46.14:12345
#end
L 01/05/2025 - 22:56:06: rcon from "123.253.249.142:12345": command "status"
`

	want := []PlayerMetric{
		{
			UserID:    10,
			Name:      "Негросучка",
			SteamID:   "1037635443",
			Connected: "32:53",
			Ping:      71,
			Loss:      0,
			State:     `active`,
			Rate:      80000,
			Adr:       `123.31.25.81:12345`,
		},
		{
			UserID:    11,
			Name:      "фея технологий",
			SteamID:   `1175148404`,
			Connected: "32:43",
			Ping:      93,
			Loss:      0,
			State:     `active`,
			Rate:      80000,
			Adr:       `123.219.46.14:12345`,
		},
	}
	expectParsedPlayers(t, raw, want)

}

func TestParseStatusResponse_Human(t *testing.T) {
	raw := `version : 41/41 0 secure  
steamid : [A:1:2405480470:46530] (90271839781173270)
udp/ip  :  172.18.0.3:25636 os(Linux) type(dedicated)
sourcetv:  port 25641, delay 30.0s
players : 10 humans, 1 bots (15 max) (not hibernating)
edicts : 1090 used of 2048 max
gamestate: DOTA_GAMERULES_STATE_GAME_IN_PROGRESS Times: Transition=291.61 Current=1442.17
# userid name uniqueid connected ping loss state rate adr
#  3 2 "Fosoroy через атос" [U:1:126880756] 28:53 87 0 active 80000 123.214.4.218:12345`

	want := []PlayerMetric{
		{
			UserID:    2,
			Name:      "Fosoroy через атос",
			SteamID:   "126880756",
			Connected: "28:53",
			Ping:      87,
			Loss:      0,
			State:     "active",
			Rate:      80000,
			Adr:       "123.214.4.218:12345",
		},
	}
	expectParsedPlayers(t, raw, want)
}

func TestParseStatusRow_InvalidRow(t *testing.T) {
	row := `# 4 "BotPlayer"`
	if m := parseStatusRow(row); m != nil {
		t.Errorf("expected nil for invalid/bot row, got %+v", m)
	}
}

func expectParsedPlayers(t *testing.T, raw string, want []PlayerMetric) {

	got := ParseStatusResponse(raw)

	if len(got) != len(want) {
		t.Fatalf("expected %d players, got %d", len(want), len(got))
	}

	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Errorf("player %d mismatch:\n got  %+v\n want %+v", i, got[i], want[i])
		}
	}
}
