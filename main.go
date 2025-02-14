package main

import (
	"log"
	"net/http"
	"os"

	"github.com/newnes25/helloworld/games"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// WebSocket Odd/Even game route
	http.HandleFunc("/ws/odd_even", games.HandleOddEvenGame)

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
