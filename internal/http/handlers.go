package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sidecar/internal/mapper"
	"sidecar/internal/models"
	"sidecar/internal/rabbit"
	"sidecar/internal/redis"
	"sidecar/internal/srcds/log_parser"
	"sidecar/internal/srcds/metrics"
	"sidecar/internal/state"
	"sidecar/internal/util"
	"strconv"
	"time"

	d2cmodels "github.com/dota2classic/d2c-go-models/models"
)

func HandleJSONPost[T any](handler func(T, http.ResponseWriter)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data T
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Printf("Error decoding body: %s", err)
			http.Error(w, fmt.Sprintf("Failed to parse JSON: %v", err), http.StatusBadRequest)
			return
		}

		handler(data, w)
	}
}

func HandleLiveMatch(data models.LiveMatchDto, w http.ResponseWriter) {
	mapped := mapper.MapLiveMatchUpdatedEvent(data)
	redis.PublishLiveMatch(&mapped)
	w.WriteHeader(http.StatusOK)
}

func HandleMatchFailed(d models.MatchFailedOnSRCDS, w http.ResponseWriter) {
	log.Printf("Received failed_match: %+v", d)
	var failedPlayers []models.FailedPlayerInfo
	for _, p := range d.Players {
		if p.Connection == d2cmodels.DOTA_CONNECTION_STATE_FAILED {
			failedPlayers = append(failedPlayers, p)
		}
	}

	// Identify "good" parties (those with no failed players)
	var goodParties []string
partyLoop:
	for _, p := range d.Players {
		for _, failed := range failedPlayers {
			if p.PartyID == failed.PartyID {
				continue partyLoop
			}
		}
		goodParties = append(goodParties, p.PartyID)
	}

	// If any players failed, emit event
	if len(failedPlayers) > 0 {
		var failedIDs []string
		for _, p := range failedPlayers {
			failedIDs = append(failedIDs, strconv.FormatInt(p.SteamID, 10))
		}

		event := d2cmodels.MatchFailedEvent{
			MatchID:       d.MatchID,
			Server:        state.GlobalMatchInfo.ServerAddress,
			FailedPlayers: failedIDs,
			GoodParties:   goodParties,
		}
		rabbit.PublishMatchFailedEvent(&event)
	}
	w.WriteHeader(http.StatusOK)
}

func HandlePlayerNotLoaded(data models.PlayerNotLoadedOnSRCDS, w http.ResponseWriter) {
	log.Printf("Received player_not_loaded: %+v", data)
	event := d2cmodels.MatchFailedEvent{
		MatchID:       data.MatchID,
		Server:        data.Server,
		FailedPlayers: []string{data.SteamID},
		GoodParties:   []string{},
	}
	rabbit.PublishMatchFailedEvent(&event)
	w.WriteHeader(http.StatusOK)
}

func HandlePlayerAbandon(data models.PlayerAbandonOnSRCDS, w http.ResponseWriter) {
	log.Printf("Received player_abandon: %+v", data)
	event := d2cmodels.PlayerAbandonedEvent{
		PlayerId:     d2cmodels.PlayerIdType{Value: strconv.FormatInt(data.SteamID, 10)},
		MatchId:      data.MatchID,
		AbandonIndex: data.AbandonIndex,
		Mode:         data.Mode,
		GameState:    data.GameState,
	}
	rabbit.PublishPlayerAbandonEvent(&event)
	w.WriteHeader(http.StatusOK)
}

func HandlePlayerConnect(data models.PlayerConnectedOnSRCDS, w http.ResponseWriter) {
	log.Printf("Received player_connected: %+v", data)
	event := d2cmodels.PlayerConnectedEvent{
		PlayerId:  d2cmodels.PlayerIdType{Value: strconv.FormatInt(data.SteamID, 10)},
		MatchId:   data.MatchID,
		ServerUrl: data.Server,
		Ip:        data.IP,
	}
	redis.PublishPlayerConnectedEvent(&event)
	metrics.ObserveLoadingTime(data.LobbyType, data.Duration)
	w.WriteHeader(http.StatusOK)
}

func HandleMatchResults(data models.MatchFinishedOnSRCDS, w http.ResponseWriter) {
	log.Printf("Received match_results: %+v", data)
	time.Sleep(5 * time.Second)

	event := mapper.MapGameResults(data)
	err := log_parser.FillAdditionalDataFromLog(&event, util.GetLogFilePath())
	if err != nil {
		log.Printf("Failed to fill additional data: %v", err)
	}

	rabbit.PublishGameResultsEvent(&event)
	redis.ServerStatus(false)

	w.WriteHeader(http.StatusOK)
}
