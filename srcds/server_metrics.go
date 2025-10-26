package srcds

import (
	"fmt"
	"log"
	"os"
	"sidecar/state"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	registry     = prometheus.NewRegistry()
	pushgateway  = os.Getenv("PUSHGATEWAY_URL")
	jobName      = "gameserver_sidecar"
	registerOnce sync.Once

	// Metrics definitions
	LoadingTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "d2c_game_server_loading_time",
			Help:    "Loading time into game",
			Buckets: prometheus.LinearBuckets(15, 15, 10), // 15, 30, ... 150 seconds
		},
		[]string{"host", "lobby_type"},
	)

	CpuGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_cpu",
			Help: "CPU usage of SRCDS process",
		},
		[]string{"host", "match_id", "server_url", "lobby_type"},
	)

	FpsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_fps",
			Help: "Frames per second of SRCDS process",
		},
		[]string{"host", "match_id", "server_url", "lobby_type"},
	)

	InGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_net_in",
			Help: "Inbound network usage (bytes/sec)",
		},
		[]string{"host", "match_id", "server_url", "lobby_type"},
	)

	OutGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_net_out",
			Help: "Outbound network usage (bytes/sec)",
		},
		[]string{"host", "match_id", "server_url", "lobby_type"},
	)

	PingGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_ping",
			Help: "Ping per player",
		},
		[]string{"host", "match_id", "server_url", "lobby_type", "steam_id"},
	)

	LossGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_loss",
			Help: "Packet loss per player",
		},
		[]string{"host", "match_id", "server_url", "lobby_type", "steam_id"},
	)

	PlayerCountGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_player_count",
			Help: "Number of active players on the server",
		},
		[]string{"host", "match_id", "server_url", "lobby_type"},
	)
)

// SrcdsServerMetrics represents server stats
type SrcdsServerMetrics struct {
	CPU     float64
	In      float64
	Out     float64
	Uptime  float64
	Users   int
	FPS     float64
	Players int
}

// ParseAndRecordSrcdsMetrics parses the raw stats string into SrcdsServerMetrics
func ParseAndRecordSrcdsMetrics(statsRaw string) {
	stats, err := parseRawStats(statsRaw)
	if err != nil {
		log.Println("Error parsing stats: ", err)
		return
	}

	labels := []string{
		state.GlobalMatchInfo.Host,
		strconv.FormatInt(state.GlobalMatchInfo.MatchID, 10),
		state.GlobalMatchInfo.ServerAddress,
		string(rune(state.GlobalMatchInfo.LobbyType)),
	}

	CpuGauge.WithLabelValues(labels...).Set(stats.CPU)
	FpsGauge.WithLabelValues(labels...).Set(stats.FPS)
	InGauge.WithLabelValues(labels...).Set(stats.In)
	OutGauge.WithLabelValues(labels...).Set(stats.Out)
	PlayerCountGauge.WithLabelValues(labels...).Set(float64(stats.Players))

}

func parseRawStats(statsRaw string) (*SrcdsServerMetrics, error) {
	lines := strings.Split(statsRaw, "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("invalid stats format")
	}

	// remove first and last line
	lines = lines[1 : len(lines)-1]
	fields := strings.Fields(lines[0])

	if len(fields) < 7 {
		return nil, fmt.Errorf("not enough fields in stats line")
	}

	cpu, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil, err
	}
	in, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}
	out, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, err
	}
	uptime, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, err
	}
	users, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, err
	}
	fps, err := strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return nil, err
	}
	players, err := strconv.Atoi(fields[6])
	if err != nil {
		return nil, err
	}

	return &SrcdsServerMetrics{
		CPU:     cpu,
		In:      in,
		Out:     out,
		Uptime:  uptime,
		Users:   users,
		FPS:     fps,
		Players: players,
	}, nil
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

// Push all metrics to the Pushgateway
func PushMetrics(groupLabels map[string]string) {
	initMetrics()

	pusher := push.New(pushgateway, jobName).Gatherer(registry)
	for k, v := range groupLabels {
		pusher = pusher.Grouping(k, v)
	}

	if err := pusher.Push(); err != nil {
		log.Printf("Failed to push metrics: %v", err)
	}
}
