package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if PORT is not set
	}

	http.HandleFunc("/", handler)

	log.Printf("Server starting on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
