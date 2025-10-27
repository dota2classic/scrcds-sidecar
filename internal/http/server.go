package http

import (
	"fmt"
	"log"
	"net/http"
)

func Listen(port int) {
	log.Printf("Starting server on port %d...\n", port)
	initRoutes()
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))

}
