package http

import "net/http"

func initRoutes() {

	http.HandleFunc("/live_match", HandleJSONPost(HandleLiveMatch))

	http.HandleFunc("/failed_match", HandleJSONPost(HandleMatchFailed))

	http.HandleFunc("/player_not_loaded", HandleJSONPost(HandlePlayerNotLoaded))

	http.HandleFunc("/player_abandon", HandleJSONPost(HandlePlayerAbandon))

	http.HandleFunc("/player_connect", HandleJSONPost(HandlePlayerConnect))

	http.HandleFunc("/match_results", HandleJSONPost(HandleMatchFinished))
}
