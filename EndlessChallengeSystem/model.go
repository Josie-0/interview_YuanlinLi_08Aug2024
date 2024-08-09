package main

import (
	"sync"
	"time"
)

type Challenge struct {
	ID        string    `json:"id"`
	PlayerID  string    `json:"player_id"`
	Amount    float64   `json:"amount"`
	Won       bool      `json:"won"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

var (
	challenges          = make(map[string]Challenge)
	playerLastChallenge = make(map[string]time.Time)
	challengeMutex      sync.Mutex
	playerMutex         sync.Mutex
)
