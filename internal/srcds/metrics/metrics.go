package metrics

import (
	"log"
	"os"
	"sidecar/internal/state"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	registry     = prometheus.NewRegistry()
	pushgateway  = os.Getenv("PUSHGATEWAY_URL")
	jobName      = "gameserver_sidecar"
	registerOnce sync.Once
)

func CollectMetrics() error {
	if err := collectServerMetrics(); err != nil {
		return err
	}
	collectPlayerMetrics()
	if os.Getenv("FRAME_METRICS_ENABLED") == "true" {
		collectFrameMetrics()
	}
	pushMetrics(getGroupLabels())
	return nil
}

// Delete will delete all pushed metrics so that they don't continue indefinitely
func Delete() {
	pusher := push.New(pushgateway, jobName).Gatherer(nil) // nil because we’re deleting
	for k, v := range getGroupLabels() {
		pusher = pusher.Grouping(k, v)
	}

	if err := pusher.Delete(); err != nil {
		log.Printf("Failed to delete metrics: %v", err)
	}
}

func getGroupLabels() map[string]string {
	return map[string]string{
		"jobName": "srcds-sidecar", "host": state.GlobalMatchInfo.Host,
		"match_id": strconv.FormatInt(state.GlobalMatchInfo.MatchID, 10),
		"pod_name": os.Getenv("POD_NAME"),
	}
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
		os.Getenv("POD_NAME"),
		strconv.FormatInt(state.GlobalMatchInfo.MatchID, 10),
		state.GlobalMatchInfo.ServerAddress,
		strconv.Itoa(int(state.GlobalMatchInfo.LobbyType)),
	}
	return labels
}

// Register all metrics once
func initMetrics() {
	registerOnce.Do(func() {
		registry.MustRegister(

			// Server metrics
			CpuGauge,
			FpsGauge,
			InGauge,
			OutGauge,
			PlayerCountGauge,
			FrameWindowGauge, FrameMeanGauge, FrameP95Gauge, FrameP99Gauge,
			FrameMaxGauge, FrameOverBudgetGauge,

			// Player metrics
			PingGauge,
			LossGauge,

			// Other metrics
			LoadingTime,
		)
	})
}
