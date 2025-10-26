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

	http.HandleFunc("/live_match", handleLiveMatch)

	port := 7777
	log.Printf("Starting server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleLiveMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data models.LiveMatchDto
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Received POST request: %+v\n", data)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello %d, age %d", data.MatchID, data.Duration)))
}
