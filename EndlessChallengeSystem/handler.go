package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func ParticipateChallenge(w http.ResponseWriter, r *http.Request) {
	var request struct {
		PlayerID string  `json:"player_id"`
		Amount   float64 `json:"amount"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Amount != 20.01 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	playerMutex.Lock()
	lastChallengeTime, exists := playerLastChallenge[request.PlayerID]
	playerMutex.Unlock()

	if exists && time.Since(lastChallengeTime) < time.Minute {
		http.Error(w, "You can only participate once per minute", http.StatusTooManyRequests)
		return
	}

	challengeID := uuid.New().String()

	challenge := Challenge{
		ID:        challengeID,
		PlayerID:  request.PlayerID,
		Amount:    request.Amount,
		Timestamp: time.Now(),
		Status:    "processing",
	}

	playerMutex.Lock()
	playerLastChallenge[request.PlayerID] = time.Now()
	playerMutex.Unlock()

	challengeMutex.Lock()
	challenges[challengeID] = challenge
	challengeMutex.Unlock()

	go func() {
		time.Sleep(30 * time.Second)
		challengeMutex.Lock()
		challenge.Won = simulateChallengeOutcome()
		challenge.Status = "completed"
		challenges[challengeID] = challenge
		challengeMutex.Unlock()
	}()

	response := map[string]string{
		"challenge_id": challengeID,
		"status":       "processing",
		"message":      "Challenge is being processed. Check back later for the result.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAllChallenges(w http.ResponseWriter, r *http.Request) {
	challengeMutex.Lock()
	response := make([]Challenge, 0, len(challenges))
	for _, challenge := range challenges {
		response = append(response, challenge)
	}
	challengeMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetChallengeResult(w http.ResponseWriter, r *http.Request) {
	challengeID := mux.Vars(r)["id"]

	challengeMutex.Lock()
	challenge, exists := challenges[challengeID]
	challengeMutex.Unlock()

	if !exists {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenge)
}

func simulateChallengeOutcome() bool {
	return rand.Float64() <= 0.01
}
