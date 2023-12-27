package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for simplicity; you may want to implement origin checking in production.
		return true
	},
}

var connections = make(map[*websocket.Conn]struct{})

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called websocket endpoint")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected")
	// Add the new connection to the list
	connections[conn] = struct{}{}

	// Handle WebSocket events
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		var progress Progress
		err = json.Unmarshal(p, &progress)
		if err == nil {
			fmt.Println(progress)
			players.Progress[progress.PlayerID] = progress.Idx
			EventStruct := Event{
				Event: "players",
				Data:  players,
			}
			message, _ := json.Marshal(EventStruct)
			broadcastMessage(websocket.TextMessage, message)
		}

		if string(p) == "join" {
			fmt.Println("new player joined")
			playerID := updatePlayers()

			EventStruct := Event{
				Event: "joined",
				Data:  playerID,
			}
			message, _ := json.Marshal(EventStruct)
			conn.WriteMessage(websocket.TextMessage, message)

			EventStruct = Event{
				Event: "players",
				Data:  players,
			}
			message, _ = json.Marshal(EventStruct)
			broadcastMessage(websocket.TextMessage, message)
		}

		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
	}
}

func broadcastMessage(messageType int, message []byte) {
	// Iterate over all connections and send the message
	fmt.Println("Broadcasting ", len(connections))
	for conn := range connections {
		err := conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("Error writing message:", err)
			// Optionally handle errors or remove the connection from the list
		}
	}
}

func handleJoinLobby(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	data := updatePlayers()
	fmt.Fprintf(w, "%d", data)
}

func handleGetPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	n := getPlayers()
	fmt.Fprintf(w, "%d", n)
}

// create a user and save the number of races he completed.
func raceCompletedhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	updateRaceCompleted("anshul")
	fmt.Fprintf(w, "Hello from golang webserver")
}

func raceCnt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	races := getRaceCompleted("anshul")
	fmt.Fprintf(w, "%d", races)
}

func handleFrontPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprintf(w, "Welcome to typing race backend")
}

func main() {
	players = Players{
		N:        0,
		Progress: make(map[int]int),
	}
	createMongoClient()
	// createUser("anshul")
	// insertUserProfile("anshul")
	// updateRaceCompleted("anshul")

	fmt.Println("Hello World")
	http.HandleFunc("/raceCompleted", raceCompletedhandler)
	http.HandleFunc("/races", raceCnt)
	http.HandleFunc("/joinLobby", handleJoinLobby)
	http.HandleFunc("/players", handleGetPlayers)
	// http.HandleFunc("/", handleFrontPage)

	http.HandleFunc("/ws", handleWebSocket)
	http.ListenAndServe("10.107.107.107:8081", nil)

	// Setup signal handling to catch SIGTERM
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-stopChan
		fmt.Printf("Received signal %v. Shutting down...\n", sig)

		// Disconnect MongoDB client before exiting
		DisconnectMongoClient()

		os.Exit(0)
	}()

}
