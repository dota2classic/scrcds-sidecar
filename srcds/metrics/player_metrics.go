package metrics

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorcon/rcon"
	"github.com/prometheus/client_golang/prometheus"
)

// PlayerMetric represents a player entry from `status` command output.
type PlayerMetric struct {
	UserID    int
	Name      string
	SteamID   string
	Connected string
	Ping      int
	Loss      int
	State     string
	Rate      int
	Adr       string
}

var (
	PingGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_ping",
			Help: "Ping per player",
		},
		[]string{"match_id", "server_url", "lobby_type", "steam_id"},
	)

	LossGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "srcds_metrics_loss",
			Help: "Packet loss per player",
		},
		[]string{"match_id", "server_url", "lobby_type", "steam_id"},
	)
)

func collectPlayerMetrics(conn *rcon.Conn) {
	status, err := conn.Execute("status")
	if err != nil {
		log.Printf("Failed to execute RCON command: %v", err)
		return
	}
	parseAndRecordPlayerMetrics(status)
}

func parseAndRecordPlayerMetrics(statusRaw string) {
	stats, err := parseRawRconStatsResponse(statusRaw)
	if err != nil {
		log.Println("Error parsing status: ", err)
		return
	}

	labels := getMetricLabels()

	PingGauge.WithLabelValues(labels...).Set(stats.Out)
	LossGauge.WithLabelValues(labels...).Set(float64(stats.Players))

}

// parseStatusRow parses a single row of the `status` output.
// Returns nil if it's not a valid player row (e.g. a bot).
func parseStatusRow(row string) *PlayerMetric {
	sep := regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)
	spaced := sep.FindAllString(row, -1)
	if len(spaced) < 7 {
		// likely a bot or malformed row
		return nil
	}

	// userid is the 3rd column (index 2)
	userID, err := strconv.Atoi(spaced[2])
	if err != nil {
		return nil
	}

	usernameStart := strings.Index(row, `"`)
	usernameEnd := strings.LastIndex(row, `"`)
	if usernameStart == -1 || usernameEnd == -1 || usernameEnd <= usernameStart {
		return nil
	}

	username := strings.TrimSpace(row[usernameStart+1 : usernameEnd])
	afterUsername := row[usernameEnd:]
	spaced = sep.FindAllString(afterUsername, -1)

	if len(spaced) < 7 {
		return nil
	}

	ping, _ := strconv.Atoi(spaced[2])
	loss, _ := strconv.Atoi(spaced[3])
	rate, _ := strconv.Atoi(spaced[5])

	return &PlayerMetric{
		UserID:    userID,
		Name:      username,
		SteamID:   strings.TrimSuffix(strings.TrimPrefix(spaced[0], "[U:1:"), `]`),
		Connected: spaced[1],
		Ping:      ping,
		Loss:      loss,
		State:     spaced[4],
		Rate:      rate,
		Adr:       spaced[6],
	}
}

// parseStatusResponse parses the full `status` command output.
func parseStatusResponse(raw string) []PlayerMetric {
	lines := strings.Split(raw, "\n")
	startIdx := -1

	for i, line := range lines {
		if strings.TrimSpace(line) == "# userid name uniqueid connected ping loss state rate adr" {
			startIdx = i + 1
			break
		}
	}
	if startIdx == -1 {
		return nil
	}

	var results []PlayerMetric
	for _, line := range lines[startIdx:] {
		if strings.TrimSpace(line) == "#end" {
			break
		}
		if metric := parseStatusRow(line); metric != nil {
			results = append(results, *metric)
		}
	}
	return results
}
