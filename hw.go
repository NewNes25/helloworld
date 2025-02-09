package main

import (
	"fmt"
	"log"
	"net/http"

	"sync"

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
var mutex sync.Mutex

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	player := &Player{conn: conn}
	mutex.Lock()
	players = append(players, player)
	playerID := len(players)
	mutex.Unlock()

	log.Printf("Player %d connected\n", playerID)

	// Wait for the player to send a number (1 or 2)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		removePlayer(conn)
		return
	}

	var num int
	if msg[0] == '1' {
		num = 1
	} else if msg[0] == '2' {
		num = 2
	} else {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid input! Send '1' or '2'."))
		removePlayer(conn)
		return
	}

	player.number = num
	log.Printf("Player %d chose: %d\n", playerID, num)

	// Wait until two players are connected
	mutex.Lock()
	if len(players) == 2 {
		determineWinner()
	}
	mutex.Unlock()
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

	// Send results to both players
	p1.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("You chose %d, Opponent chose %d. %s", p1.number, p2.number, result)))
	p2.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("You chose %d, Opponent chose %d. %s", p2.number, p1.number, result)))

	// Reset for the next round
	players = []*Player{}
}

func removePlayer(conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, p := range players {
		if p.conn == conn {
			players = append(players[:i], players[i+1:]...)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	port := "8080"
	log.Println("WebSocket server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
