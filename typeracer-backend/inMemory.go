package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

var availSession AvailSession // available session to join
var sessions map[int]*Lobby   // session id -> games

type AvailSession struct {
	mu sync.Mutex
	ID int
}

type Lobby struct {
	mu          sync.Mutex
	N           int               `json:"N"`
	SessionID   int               `json:"sessionID"`
	Connections []*websocket.Conn `json:"connections"` // session id -> websocket connections
	Progress    map[int]int       `json:"progress"`    // player id -> progress
	Rank        map[int]int       `json:"rank"`        // player id -> rank
	NxtRank     int               `json:"nxtrank"`
	Timer       int               `json:"timer"`
	Open        int               `json:"open"`
}

type Event struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type Joined struct {
	PlayerID  int `json:"playerID"`
	SessionID int `json:"sessionID"`
}

type Progress struct {
	PlayerID   int `json:"playerID"`
	SessionID  int `json:"sessionID"`
	Percentage int `json:"percentage"`
}
