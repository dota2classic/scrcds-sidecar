package metrics

import (
	"log"
	"os"
	"sidecar/state"
	"strconv"
	"sync"

	"github.com/gorcon/rcon"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	registry     = prometheus.NewRegistry()
	pushgateway  = os.Getenv("PUSHGATEWAY_URL")
	jobName      = "gameserver_sidecar"
	registerOnce sync.Once
)

func CollectMetrics(conn *rcon.Conn) {
	collectServerMetrics(conn)
	collectPlayerMetrics(conn)
	pushMetrics(map[string]string{jobName: "srcds-sidecar"})
}

// PushMetrics Push all metrics to the Pushgateway
func pushMetrics(groupLabels map[string]string) {
	initMetrics()

	pusher := push.New(pushgateway, jobName).Gatherer(registry)
	for k, v := range groupLabels {
		pusher = pusher.Grouping(k, v)
	}

	if err := pusher.Push(); err != nil {
		log.Printf("Failed to push metrics: %v", err)
	}
}

func getMetricLabels() []string {

	labels := []string{
		state.GlobalMatchInfo.Host,
		strconv.FormatInt(state.GlobalMatchInfo.MatchID, 10),
		state.GlobalMatchInfo.ServerAddress,
		string(rune(state.GlobalMatchInfo.LobbyType)),
	}
	return labels
}

// Register all metrics once
func initMetrics() {
	registerOnce.Do(func() {
		registry.MustRegister(
			LoadingTime,
			CpuGauge,
			FpsGauge,
			InGauge,
			OutGauge,
			PingGauge,
			LossGauge,
			PlayerCountGauge,
		)
	})
}
