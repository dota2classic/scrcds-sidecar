package metrics

import (
	"fmt"
	"math"
	"sidecar/internal/srcds/rcon"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func frameGauge(name, help string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: help}, []string{"pod_name", "match_id", "server_url", "lobby_type"})
}

var (
	FrameWindowGauge     = frameGauge("srcds_frame_window_frames", "Frames observed in the rolling 5 second window")
	FrameMeanGauge       = frameGauge("srcds_frame_interval_mean_seconds", "Mean wall-clock interval between OnGameFrame callbacks over 5 seconds; includes normal tick pacing")
	FrameP95Gauge        = frameGauge("srcds_frame_interval_p95_seconds", "95th percentile frame interval over 5 seconds")
	FrameP99Gauge        = frameGauge("srcds_frame_interval_p99_seconds", "99th percentile frame interval over 5 seconds")
	FrameMaxGauge        = frameGauge("srcds_frame_interval_max_seconds", "Maximum frame interval over 5 seconds")
	FrameOverBudgetGauge = frameGauge("srcds_frame_over_budget_ratio", "Fraction of frame intervals exceeding 1.5 times the configured tick interval over 5 seconds")
)

func parseFrameMetrics(raw string) (map[string]float64, error) {
	start := strings.Index(raw, "D2C_FRAME ")
	if start < 0 {
		return nil, fmt.Errorf("frame telemetry plugin unavailable")
	}
	line := strings.SplitN(raw[start:], "\n", 2)[0]
	values := make(map[string]float64)
	for _, field := range strings.Fields(line)[1:] {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			continue
		}
		number, err := strconv.ParseFloat(parts[1], 64)
		if err != nil || math.IsNaN(number) || math.IsInf(number, 0) || number < 0 {
			return nil, fmt.Errorf("invalid frame metric")
		}
		values[parts[0]] = number
	}
	for _, key := range []string{"frames", "mean_ms", "p95_ms", "p99_ms", "max_ms", "over_ratio"} {
		if _, ok := values[key]; !ok {
			return nil, fmt.Errorf("missing frame metric %s", key)
		}
	}
	if values["over_ratio"] > 1 {
		return nil, fmt.Errorf("invalid frame ratio")
	}
	return values, nil
}

func collectFrameMetrics() {
	raw, err := rcon.RunRconCommand("d2c_frame_stats")
	if err != nil {
		return
	}
	values, err := parseFrameMetrics(raw)
	if err != nil {
		return
	}
	labels := getMetricLabels()
	FrameWindowGauge.WithLabelValues(labels...).Set(values["frames"])
	FrameMeanGauge.WithLabelValues(labels...).Set(values["mean_ms"] / 1000)
	FrameP95Gauge.WithLabelValues(labels...).Set(values["p95_ms"] / 1000)
	FrameP99Gauge.WithLabelValues(labels...).Set(values["p99_ms"] / 1000)
	FrameMaxGauge.WithLabelValues(labels...).Set(values["max_ms"] / 1000)
	FrameOverBudgetGauge.WithLabelValues(labels...).Set(values["over_ratio"])
}
