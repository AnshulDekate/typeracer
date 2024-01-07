package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for simplicity; you may want to implement origin checking in production.
		return true
	},
}

func countDown(sessionID int) {
	lobby := sessions[sessionID]
	total := 20 // 20 seconds countdown
	for i := 0; i <= total; i++ {
		lobby.mu.Lock()
		lobby.Timer = i
		if i == 5 {
			lobby.Open = 0
		}
		sessions[sessionID] = lobby

		// send players event
		EventStruct := Event{
			Event: "players",
			Data:  lobby,
		}
		message, _ := json.Marshal(EventStruct)
		broadcastMessage(websocket.TextMessage, message, lobby.Connections)
		lobby.mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func launchBot(sessionID int, seed int) {
	time.Sleep(10 * time.Second)
	lobby := sessions[sessionID]
	lobby.mu.Lock()
	lobby.N = lobby.N + 1
	lobby.Progress[lobby.N] = 0
	playerID := lobby.N
	sendPlayersMessage(lobby)
	lobby.mu.Unlock()

	idx := 0

	for {

		lobby.mu.Lock()
		if lobby.Timer >= 20 {
			idx++
			percentage := (idx * 100) / 23
			lobby.Progress[playerID] = percentage
			if percentage == 100 {
				lobby.Rank[playerID] = lobby.NxtRank
				lobby.NxtRank = lobby.NxtRank + 1
				sendPlayersMessage(lobby)
				lobby.mu.Unlock()
				break
			}
			sendPlayersMessage(lobby)
		}
		lobby.mu.Unlock()

		rand.Seed(int64(seed) + time.Now().UnixNano())
		speed := rand.Intn(seed) + 30
		time.Sleep(time.Minute / time.Duration(speed))
	}

}

func CreateNewLobby() {
	fmt.Println("creating new lobby")
	availSession.ID++

	lobby := Lobby{
		SessionID:   availSession.ID,
		N:           1,
		Progress:    make(map[int]int),
		Rank:        make(map[int]int),
		Connections: make([]*websocket.Conn, 0),
		NxtRank:     1,
		Open:        1,
	}
	lobby.Progress[1] = 0

	sessions[availSession.ID] = &lobby

}

func sendJoinedMessage(conn *websocket.Conn, lobby *Lobby) {
	// send joined event containing playerid and session id
	fmt.Println("sending joined event")
	joined := Joined{
		PlayerID:  lobby.N,
		SessionID: lobby.SessionID,
	}
	EventStruct := Event{
		Event: "joined",
		Data:  joined,
	}
	message, _ := json.Marshal(EventStruct)
	conn.WriteMessage(websocket.TextMessage, message)
}

func sendPlayersMessage(lobby *Lobby) {
	// send players event, update the shared data of the session
	fmt.Println("sending players event")
	EventStruct := Event{
		Event: "players",
		Data:  lobby,
	}
	message, _ := json.Marshal(EventStruct)

	broadcastMessage(websocket.TextMessage, message, lobby.Connections)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Create new websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected")

	// Handle WebSocket events
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("got websocket message", string(p))

		if string(p) == "join" {
			fmt.Println("joining")
			// join should be synchronized across websockets
			// create new session or join existing

			availSession.mu.Lock()

			lobby := sessions[availSession.ID] //reference

			lobby.mu.Lock()

			joined := 0
			if lobby.Open == 1 {
				joined = 1
				fmt.Println("joining existing lobby")
				// only 5 people at max in lobby, before last 5 second in timer close the lobby
				lobby.N = lobby.N + 1
				lobby.Progress[lobby.N] = 0

				if lobby.N >= 7 {
					lobby.Open = 0
				}
				sessions[availSession.ID] = lobby

				fmt.Println("adding new connection to list")
				lobby.Connections = append(lobby.Connections, conn)

				sendJoinedMessage(conn, lobby)
				sendPlayersMessage(lobby)
			}

			lobby.mu.Unlock()

			if joined == 0 {
				fmt.Println("everything full, creating new lobby")
				CreateNewLobby()
				newlobby := sessions[availSession.ID]
				fmt.Println("adding new connection to list")
				newlobby.Connections = append(newlobby.Connections, conn)
				newlobby.mu.Lock()
				sendJoinedMessage(conn, newlobby)
				sendPlayersMessage(lobby)
				newlobby.mu.Unlock()
				go countDown(availSession.ID)
				go launchBot(availSession.ID, 5)
				go launchBot(availSession.ID, 20)
			}

			availSession.mu.Unlock()
		} else {
			fmt.Println("progressing")
			var progress Progress
			err = json.Unmarshal(p, &progress)
			if err == nil {
				fmt.Println("event progress", progress)

				lobby := sessions[progress.SessionID]
				lobby.mu.Lock()

				lobby.Progress[progress.PlayerID] = progress.Percentage

				if progress.Percentage == 100 {
					lobby.Rank[progress.PlayerID] = lobby.NxtRank
					lobby.NxtRank = lobby.NxtRank + 1
				}

				sessions[progress.SessionID] = lobby

				EventStruct := Event{
					Event: "players",
					Data:  lobby,
				}
				message, _ := json.Marshal(EventStruct)
				broadcastMessage(websocket.TextMessage, message, lobby.Connections)
				lobby.mu.Unlock()
			} else {
				fmt.Println(err)
			}
		}
	}
}

func broadcastMessage(messageType int, message []byte, sessionConnections []*websocket.Conn) {
	// Iterate over all connections and send the message
	fmt.Println("Broadcasting ", len(sessionConnections))
	for _, conn := range sessionConnections {
		err := conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("Error writing message:", err)
			// Optionally handle errors or remove the connection from the list
		}
	}
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
	fmt.Println("calling get races")
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
	// createMongoClient()
	// init global variables
	availSession = AvailSession{ID: 0}
	sessions = make(map[int]*Lobby)
	sessions[availSession.ID] = &Lobby{Open: 0}

	// createUser("anshul")
	// insertUserProfile("anshul")
	// updateRaceCompleted("anshul")

	// Routes
	// http.HandleFunc("/raceCompleted", raceCompletedhandler)
	// http.HandleFunc("/races", raceCnt)
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
