package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	ConnID   *websocket.Conn
	Username string
	Number   string
}

type Lobby struct {
	Code      string
	Players   []*Player
	CreatedAt int64
}

var lobbies = make(map[string]Lobby)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		hostParts := strings.Split(r.Host, ":")
		ipAddress := hostParts[0]

		allowedOrigins := map[string]bool{
			"127.0.0.1": true,
			"localhost": true,
		}
		return allowedOrigins[ipAddress]
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Handle WebSocket messages
	for {
		// Read message from client
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		// Echo message back to client
		if err := conn.WriteMessage(messageType, []byte("{\"message\":\"Message recived!\"}")); err != nil {
			log.Println("Write error:", err)
			break
		}

		HandleMessage(string(message), conn)

	}
}

func HandleMessage(message string, conn *websocket.Conn) {

	var msg map[string]interface{}
	json.Unmarshal([]byte(message), &msg)

	action := msg["action"].(string)
	switch action {
	case "create_lobby":
		handleCreateLobby(msg["username"].(string), conn)
	case "join_lobby":
		handleJoinLobby(msg["code"].(string), msg["username"].(string), conn)
	case "start_game":
		handleStartGame(msg["code"].(string))
	case "guess_number":
		handleGameLogic(msg["code"].(string), msg["guess"].(string), msg["player"].(string))
	}

}

func handleStartGame(code string) {
	// find lobby by code
	var lobby Lobby
	for _, l := range lobbies {
		if l.Code == code {
			lobby = Lobby{
				Code:      code,
				Players:   l.Players,
				CreatedAt: l.CreatedAt,
			}
		}
	}

	//notify client to render game page
	playerUsernames := make([]string, len(lobby.Players))
	for i, player := range lobby.Players {
		playerUsernames[i] = player.Username
		player.Number = generateNumber()
	}

	message := map[string]interface{}{
		"event":         "start_game",
		"player":        playerUsernames[0],
		"code":          lobby.Code,
		"player_number": lobby.Players[0].Number,
	}

	// send player one the message
	notifyPlayer(lobby, message, 0)

	message = map[string]interface{}{
		"event":         "start_game",
		"player":        playerUsernames[1],
		"code":          lobby.Code,
		"player_number": lobby.Players[1].Number,
	}

	// send player two the message
	notifyPlayer(lobby, message, 1)

	time.Sleep(2 * time.Second)

	message = map[string]interface{}{
		"event":     "run_game",
		"your_turn": true,
	}
	// send player one the message
	notifyPlayer(lobby, message, 0)
}

func generateNumber() string {
	rand.Seed(time.Now().UnixNano())

	// Shuffle the digits 0-9
	digits := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Shuffle(len(digits), func(i, j int) {
		digits[i], digits[j] = digits[j], digits[i]
	})

	// Form a 4-digit number using the shuffled digits
	fourDigitNumber := fmt.Sprintf("%d%d%d%d", digits[0], digits[1], digits[2], digits[3])

	return fourDigitNumber
}

func notifyPlayer(lobby Lobby, message map[string]interface{}, dest int) {
	// Broadcast the message to all players in the lobby
	p := lobby.Players[dest]

	err := p.ConnID.WriteJSON(message)
	if err != nil {
		log.Println("Error sending message to player:", err)
	}

}

func handleGameLogic(code string, guess string, player string) {
	var lobby Lobby
	var player_number string
	var guessing_player_idx int
	var waiting_player_idx int
	for _, l := range lobbies {
		if l.Code == code {
			lobby = Lobby{
				Code:      code,
				Players:   l.Players,
				CreatedAt: l.CreatedAt,
			}
			for idx, p := range l.Players {
				if p.Username != player {
					player_number = p.Number
					waiting_player_idx = idx
					break
				}
			}
		}
	}

	if waiting_player_idx == 1 {
		guessing_player_idx = 0
	} else {
		guessing_player_idx = 1
	}

	var numbers_guessed int32 = 0
	for idx, _ := range player_number {
		if player_number[idx] == guess[idx] {
			log.Println(player_number[idx], guess[idx])
			numbers_guessed += 1
		}
	}

	if numbers_guessed == 4 {
		message := map[string]interface{}{
			"event":     "game_finished",
			"isGameWon": true,
		}
		notifyPlayer(lobby, message, guessing_player_idx)

		message = map[string]interface{}{
			"event":     "game_finished",
			"isGameWon": false,
		}
		notifyPlayer(lobby, message, waiting_player_idx)

		// delete session data
		fmt.Println(lobbies)
		delete(lobbies, code)
		fmt.Println(lobbies)

		return
	}

	message := map[string]interface{}{
		"event":   "guess_response",
		"guessed": numbers_guessed,
	}
	// send player one the message
	notifyPlayer(lobby, message, guessing_player_idx)

	message = map[string]interface{}{
		"event":        "run_game",
		"your_turn":    true,
		"guess_number": guess,
	}
	// send player one the message
	notifyPlayer(lobby, message, waiting_player_idx)
	// send player one the message

}

func handleCreateLobby(username string, conn *websocket.Conn) {
	player := &Player{
		ConnID:   conn,
		Username: username,
	}

	code := generateCode(6)

	lobby := Lobby{
		Code:      code,
		Players:   []*Player{},
		CreatedAt: time.Now().Unix(),
	}
	lobby.Players = append(lobby.Players, player)

	lobbies[code] = lobby

	message := fmt.Sprintf("{\"event\":\"lobby_created\",\"code\":\"%s\",\"player1\":\"%s\"}", lobby.Code, player.Username)
	if err := conn.WriteMessage(1, []byte(message)); err != nil {
		log.Println("Write error:", err)
	}

}

func handleJoinLobby(code string, username string, conn *websocket.Conn) {

	// Get or create the lobby
	lobby, ok := lobbies[code]
	if !ok {
		// Create a new lobby if it doesn't exist
		fmt.Println("Lobby doesn't exists!")
		return
	}

	// Create a new player
	player := &Player{
		ConnID:   conn,
		Username: username,
	}

	// Add the player to the lobby
	lobby.Players = append(lobby.Players, player)

	lobbies[code] = lobby
	// Notify all players in the lobby about the new player
	notifyPlayerJoined(lobby)
}

func notifyPlayerJoined(lobby Lobby) {

	playerUsernames := make([]string, len(lobby.Players))
	for i, player := range lobby.Players {
		playerUsernames[i] = player.Username
	}
	message := map[string]interface{}{
		"event":  "player_joined",
		"player": playerUsernames[1],
		"code":   lobby.Code,
		"leader": playerUsernames[0],
	}

	message2 := map[string]interface{}{
		"event":  "join_lobby",
		"player": playerUsernames[1],
		"code":   lobby.Code,
		"leader": playerUsernames[0],
	}

	// Broadcast the message to all players in the lobby
	for idx, p := range lobby.Players {
		if idx == 0 {
			err := p.ConnID.WriteJSON(message)
			if err != nil {
				log.Println("Error sending message to player:", err)
			}
		} else {
			err := p.ConnID.WriteJSON(message2)
			if err != nil {
				log.Println("Error sending message to player:", err)
			}
		}
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}

func main() {
	// Define the route for WebSocket connections
	http.HandleFunc("/ws", handleConnection)

	// Start the HTTP server on port 8080
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
