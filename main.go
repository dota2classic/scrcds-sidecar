package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"./models"
)

func main() {
	go fileUploaderTask()

	http.HandleFunc("/live_match", handleLiveMatch)

	port := 7777
	log.Printf("Starting server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func fileUploaderTask() {
	dir := "./experiment" // directory to watch

	// Start ticker for periodic file listing
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for range ticker.C {
			listFiles(dir)
		}
	}()
}

func listFiles(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("Failed to list files in %s: %v", dir, err)
		return
	}

	fileNames := make([]string, 0, len(files))
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}

	log.Printf("Files in %s: %v", dir, fileNames)
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
	w.Write([]byte(fmt.Sprintf("Hello %s, age %d", data.MatchID, data.Duration)))
}
