package metrics

import (
	"strconv"

	"github.com/dota2classic/d2c-go-models/models"
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
)

func ObserveLoadingTime(lobbyType models.MatchmakingMode, duration float64) {
	LoadingTime.WithLabelValues(strconv.Itoa(int(lobbyType))).Observe(duration)
}
