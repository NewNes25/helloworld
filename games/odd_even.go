package games

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

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

func HandleOddEvenGame(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	defer conn.Close()

	player := &Player{conn: conn}
	players = append(players, player)

	// Channel to receive number
	numberChan := make(chan int)

	// Goroutine to wait for player input
	go func() {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		if msg[0] != '1' && msg[0] != '2' {
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid input! Send '1' or '2'."))
			return
		}
		numberChan <- int(msg[0] - '0')
	}()

	// Wait for input with a timeout of 10 seconds
	select {
	case num := <-numberChan:
		player.number = num
	case <-time.After(10 * time.Second):
		player.number = rand.Intn(2) + 1
		player.conn.WriteMessage(websocket.TextMessage, []byte("Timeout! You were assigned: "+fmt.Sprint(player.number)))
	}

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

	p1.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("You: %d, Opponent: %d. %s", p1.number, p2.number, result)))
	p2.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("You: %d, Opponent: %d. %s", p2.number, p1.number, result)))

	players = []*Player{}
}
