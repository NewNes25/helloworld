package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Player struct {
	conn   *websocket.Conn
	number int
}

var players []*Player

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	player := &Player{conn: conn}
	players = append(players, player)

	// Wait for number
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read error:", err)
		return
	}

	if msg[0] != '1' && msg[0] != '2' {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid input! Send '1' or '2'."))
		return
	}
	player.number = int(msg[0] - '0')

	// Check if two players are connected
	if len(players) == 2 {
		determineWinner()
	}
}

func determineWinner() {
	p1, p2 := players[0], players[1]
	sum := p1.number + p2.number
	var result string

	if sum%2 == 0 {
		result = "Player 2 wins!"
	} else {
		result = "Player 1 wins!"
	}

	// Send results to players
	p1.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("You: %d, Opponent: %d. %s", p1.number, p2.number, result)))
	p2.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("You: %d, Opponent: %d. %s", p2.number, p1.number, result)))

	// Reset players for the next round
	players = []*Player{}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve static files (client.html)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Handle WebSocket connections
	http.HandleFunc("/ws", handleConnection)

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
