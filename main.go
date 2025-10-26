package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sidecar/models"
	"sidecar/srcds"
)

func main() {
	go srcds.RunHeartbeatPoller()

	http.HandleFunc("/live_match", handleJSONPost(handleLiveMatch))

	http.HandleFunc("/failed_match", handleJSONPost(handleMatchFailed))

	http.HandleFunc("/player_not_loaded", handleJSONPost(handlePlayerNotLoaded))

	http.HandleFunc("/player_abandon", handleJSONPost(handlePlayerAbandon))

	http.HandleFunc("/player_connect", handleJSONPost(handlePlayerConnect))

	http.HandleFunc("/match_results", handleJSONPost(handleMatchFinished))

	port := 7777
	log.Printf("Starting server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
}

func handleJSONPost[T any](handler func(T, http.ResponseWriter)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data T
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse JSON: %v", err), http.StatusBadRequest)
			return
		}

		handler(data, w)
	}
}

func handleLiveMatch(data models.LiveMatchDto, w http.ResponseWriter) {
	log.Printf("Received live_match: %+v", data)
	w.WriteHeader(http.StatusOK)
}

func handleMatchFailed(data models.MatchFailedOnSRCDS, w http.ResponseWriter) {
	log.Printf("Received failed_match: %+v", data)
	w.WriteHeader(http.StatusOK)
}

func handlePlayerNotLoaded(data models.PlayerNotLoadedOnSRCDS, w http.ResponseWriter) {
	log.Printf("Received player_not_loaded: %+v", data)
	w.WriteHeader(http.StatusOK)
}

func handlePlayerAbandon(data models.PlayerAbandonOnSRCDS, w http.ResponseWriter) {
	log.Printf("Received player_abandon: %+v", data)
	w.WriteHeader(http.StatusOK)
}

func handlePlayerConnect(data models.PlayerConnectedOnSRCDS, w http.ResponseWriter) {
	log.Printf("Received player_connected: %+v", data)
	w.WriteHeader(http.StatusOK)
}

func handleMatchFinished(data models.MatchFinishedOnSRCDS, w http.ResponseWriter) {
	log.Printf("Received player_connected: %+v", data)
	w.WriteHeader(http.StatusOK)
}
