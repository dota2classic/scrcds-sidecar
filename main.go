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

	http.HandleFunc("/live_match", handleJSONPost(func(data models.LiveMatchDto, w http.ResponseWriter) {
		log.Printf("Received live_match: %+v", data)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello %d, duration %d", data.MatchID, data.Duration)
	}))

	http.HandleFunc("/failed_match", handleJSONPost(func(data models.MatchFailedOnSRCDS, w http.ResponseWriter) {
		log.Printf("Received failed_match: %+v", data)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Match %d failed", data.MatchID)
	}))

	port := 7777
	log.Printf("Starting server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
}

// generic JSON handler wrapper
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
