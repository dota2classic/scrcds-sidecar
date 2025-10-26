package srcds

import (
	"reflect"
	"testing"
)

func TestParseStatsResponse(t *testing.T) {
	raw := `CPU   In    Out   Uptime  Users   FPS    Players
31.01  2.00  4.00       3     0   59.94       9
L 01/06/2025 - 00:03:30: rcon from \"156.253.249.142:55454\": command \"stats\"
`

	want := &SrcdsServerMetrics{
		CPU:     31.01,
		FPS:     59.94,
		Uptime:  3,
		In:      2.0,
		Out:     4.0,
		Players: 9,
	}

	got, err := ParseStatsResponse(raw)
	if err != nil {
		t.Fatalf("ParseStatsResponse returned error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseStatsResponse = %+v; want %+v", got, want)
	}
}

func TestParseStatsResponse_Invalid(t *testing.T) {
	// Too few lines
	_, err := ParseStatsResponse("only one line")
	if err == nil {
		t.Errorf("expected error for invalid input, got nil")
	}

	// Too few fields
	_, err = ParseStatsResponse(`
header
1 2 3
footer
`)
	if err == nil {
		t.Errorf("expected error for insufficient fields, got nil")
	}
}
