package metrics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sidecar/internal/state"
	"strings"
	"testing"
)

func TestConcurrentMatchPushAndDeleteAreIsolated(t *testing.T) {
	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		requests = append(requests, r.Method+" "+r.URL.Path)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()
	oldGateway, oldState := pushgateway, state.GlobalMatchInfo
	defer func() { pushgateway, state.GlobalMatchInfo = oldGateway, oldState }()
	pushgateway = server.URL
	t.Setenv("POD_NAME", "game-a")
	state.GlobalMatchInfo = &state.MatchInfo{Host: "same-host", MatchID: 101}
	FpsGauge.WithLabelValues("game-a", "101", "same-host:1", "1").Set(60)
	pushMetrics(getGroupLabels())
	t.Setenv("POD_NAME", "game-b")
	state.GlobalMatchInfo = &state.MatchInfo{Host: "same-host", MatchID: 102}
	FpsGauge.WithLabelValues("game-b", "102", "same-host:2", "1").Set(30)
	pushMetrics(getGroupLabels())
	Delete()
	if len(requests) != 3 {
		t.Fatalf("unexpected requests: %v", requests)
	}
	if !strings.Contains(requests[0]+"/", "/match/101/") || !strings.Contains(requests[1]+"/", "/match/102/") || !strings.HasPrefix(requests[2], "DELETE ") || !strings.Contains(requests[2]+"/", "/match/102/") {
		t.Fatalf("push/delete crossed match groups: %v", requests)
	}
}
