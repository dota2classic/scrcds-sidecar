package metrics

import "testing"

func TestParseFrameMetrics(t *testing.T) {
	v, err := parseFrameMetrics("D2C_FRAME frames=300 mean_ms=16.7 p95_ms=18 p99_ms=40 max_ms=120 over_ratio=0.02\n")
	if err != nil || v["max_ms"] != 120 || v["frames"] != 300 {
		t.Fatalf("unexpected parse: %v %v", v, err)
	}
	for _, bad := range []string{"Unknown command", "D2C_FRAME frames=2", "D2C_FRAME frames=NaN", "D2C_FRAME frames=-1"} {
		if _, err := parseFrameMetrics(bad); err == nil {
			t.Fatalf("accepted invalid response: %s", bad)
		}
	}
}
