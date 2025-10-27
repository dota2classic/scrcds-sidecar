package metrics

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gorcon/rcon"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	LoadingTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "d2c_game_server_loading_time",
			Help:    "Loading time into game",
			Buckets: prometheus.LinearBuckets(15, 15, 10), // 15, 30, ... 150 seconds
		},
		[]string{"lobby_type"},
	)

	CpuGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_cpu",
			Help: "CPU usage of SRCDS process",
		},
		[]string{"match_id", "server_url", "lobby_type"},
	)

	FpsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_fps",
			Help: "Frames per second of SRCDS process",
		},
		[]string{"match_id", "server_url", "lobby_type"},
	)

	InGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_net_in",
			Help: "Inbound network usage (bytes/sec)",
		},
		[]string{"match_id", "server_url", "lobby_type"},
	)

	OutGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_net_out",
			Help: "Outbound network usage (bytes/sec)",
		},
		[]string{"match_id", "server_url", "lobby_type"},
	)
	PlayerCountGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_player_count",
			Help: "Number of active players on the server",
		},
		[]string{"match_id", "server_url", "lobby_type"},
	)
)

// ServerMetrics represents server stats
type ServerMetrics struct {
	CPU     float64
	In      float64
	Out     float64
	Uptime  float64
	Users   int
	FPS     float64
	Players int
}

func collectServerMetrics(conn *rcon.Conn) {
	stats, err := conn.Execute("stats")
	if err != nil {
		log.Printf("Failed to execute RCON command: %v", err)
		return
	}
	parseAndRecordSrcdsMetrics(stats)
}

// parseAndRecordSrcdsMetrics parses the raw stats string into ServerMetrics
func parseAndRecordSrcdsMetrics(statsRaw string) {
	stats, err := parseRawRconStatsResponse(statsRaw)
	if err != nil {
		log.Println("Error parsing stats: ", err)
		return
	}

	labels := getMetricLabels()

	CpuGauge.WithLabelValues(labels...).Set(stats.CPU)
	FpsGauge.WithLabelValues(labels...).Set(stats.FPS)
	InGauge.WithLabelValues(labels...).Set(stats.In)
	OutGauge.WithLabelValues(labels...).Set(stats.Out)
	PlayerCountGauge.WithLabelValues(labels...).Set(float64(stats.Players))

}

func parseRawRconStatsResponse(statsRaw string) (*ServerMetrics, error) {
	lines := strings.Split(statsRaw, "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("invalid stats format")
	}

	// remove first and last line
	lines = lines[1 : len(lines)-1]
	fields := strings.Fields(lines[0])

	if len(fields) < 7 {
		fmt.Printf("Fields %s", fields)
		return nil, fmt.Errorf("not enough fields in stats line: %d", len(fields))
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

	return &ServerMetrics{
		CPU:     cpu,
		In:      in,
		Out:     out,
		Uptime:  uptime,
		Users:   users,
		FPS:     fps,
		Players: players,
	}, nil
}
