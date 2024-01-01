package main

var currSessionID int = 1

var sessions map[int]Players

type Players struct {
	SessionID int         `json:"sessionID"`
	N         int         `json:"N"`
	Progress  map[int]int `json:"progress"`
	Rank      map[int]int `json:"rank"`
	NxtRank   int         `json:"nxtrank"`
	Timer     int         `json:"timer"`
}

type Event struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type Progress struct {
	SessionID  int `json:"sessionID"`
	PlayerID   int `json:"playerID"`
	Percentage int `json:"percentage"`
}

type Joined struct {
	PlayerID  int `json:"playerID"`
	SessionID int `json:"sessionID"`
}

// func updatePlayers(players Players) int {
// 	players.N = players.N + 1
// 	playerID := players.N
// 	players.Progress[playerID] = 0
// 	return playerID
// }

// func getPlayers(players Players) int {
// 	return players.N
// }
