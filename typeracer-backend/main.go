package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

type ConnList struct {
	mu  sync.Mutex
	all []*websocket.Conn
}

var connections = make(map[int]*ConnList)

func countDown(sessionID int) {
	total := 10 // 10 seconds countdown
	for i := 0; i <= total; i++ {
		players := sessions[sessionID]
		players.Timer = i
		if i == 5 {
			players.Open = 0
		}
		sessions[sessionID] = players

		// send players event
		EventStruct := Event{
			Event: "players",
			Data:  players,
		}
		message, _ := json.Marshal(EventStruct)
		connections[players.SessionID].mu.Lock()
		broadcastMessage(websocket.TextMessage, message, connections[players.SessionID].all)
		connections[players.SessionID].mu.Unlock()
		time.Sleep(1 * time.Second)
	}
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
		fmt.Println("got websocket message")

		if string(p) == "join" {
			// Player wants to join available session

			// with current session id, get players data
			players := sessions[currSessionID]

			fmt.Println("new player joined")
			// if this is the first player initialize the struct
			if players.SessionID == 0 {
				fmt.Println("creating new lobby")
				players = Players{
					SessionID: currSessionID,
					N:         1,
					Progress:  make(map[int]int),
					Rank:      make(map[int]int),
					NxtRank:   1,
					Open:      1,
				}
				players.Progress[1] = 0
				sessionID := currSessionID
				if _, ok := connections[sessionID]; !ok {
					connections[sessionID] = &ConnList{}
				}
				if connections[sessionID].all == nil {
					connections[sessionID].all = make([]*websocket.Conn, 0)
				}
				sessions[currSessionID] = players
			} else {
				fmt.Println("joining existing lobby")
				// only 5 people at max in lobby, before last 5 second in timer close the lobby
				if players.Open == 1 {
					players.N = players.N + 1
					players.Progress[players.N] = 0
					if players.N == 2 {
						go countDown(players.SessionID)
					} else if players.N == 5 {
						players.Open = 0
					}
					sessions[currSessionID] = players
					fmt.Println(players)
				} else {
					fmt.Println("everything full, creating new lobby")
					currSessionID++
					players = Players{
						SessionID: currSessionID,
						N:         1,
						Progress:  make(map[int]int),
						Rank:      make(map[int]int),
						NxtRank:   1,
						Open:      1,
					}
					players.Progress[1] = 0
					sessionID := currSessionID
					if _, ok := connections[sessionID]; !ok {
						connections[sessionID] = &ConnList{}
					}
					if connections[sessionID].all == nil {
						connections[sessionID].all = make([]*websocket.Conn, 0)
					}
					sessions[currSessionID] = players
				}
			}

			// send joined event containing playerid and session id
			fmt.Println("sending joined event")
			joined := Joined{
				PlayerID:  players.N,
				SessionID: players.SessionID,
			}
			EventStruct := Event{
				Event: "joined",
				Data:  joined,
			}
			message, _ := json.Marshal(EventStruct)
			conn.WriteMessage(websocket.TextMessage, message)

			fmt.Println("adding new connection to list")
			connections[players.SessionID].mu.Lock()
			connections[players.SessionID].all = append(connections[players.SessionID].all, conn)
			connections[players.SessionID].mu.Unlock()

			// send players event, update the shared data of the session (currently broadcasting to all sessions)
			fmt.Println("sending players event")
			EventStruct = Event{
				Event: "players",
				Data:  players,
			}
			message, _ = json.Marshal(EventStruct)

			connections[players.SessionID].mu.Lock()
			broadcastMessage(websocket.TextMessage, message, connections[players.SessionID].all)
			connections[players.SessionID].mu.Unlock()

		} else {
			var progress Progress
			err = json.Unmarshal(p, &progress)
			if err == nil {
				// with current session id, get players data
				players := sessions[progress.SessionID]

				fmt.Println("event progress", progress, players)

				players.Progress[progress.PlayerID] = progress.Percentage

				if progress.Percentage == 100 {
					players.Rank[progress.PlayerID] = players.NxtRank
					players.NxtRank = players.NxtRank + 1
				}

				sessions[progress.SessionID] = players

				EventStruct := Event{
					Event: "players",
					Data:  players,
				}
				message, _ := json.Marshal(EventStruct)
				connections[players.SessionID].mu.Lock()
				broadcastMessage(websocket.TextMessage, message, connections[players.SessionID].all)
				connections[players.SessionID].mu.Unlock()
			} else {
				fmt.Println(err)
			}
		}
	}
}

func broadcastMessage(messageType int, message []byte, sessionConnections []*websocket.Conn) {
	// Iterate over all connections and send the message
	fmt.Println("Broadcasting ", len(connections))
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
	sessions = make(map[int]Players)
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
