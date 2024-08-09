// handler_test.go
package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParticipateChallenge(t *testing.T) {
	reqBody := `{"player_id": "player123", "amount": 20.01}`
	req := httptest.NewRequest(http.MethodPost, "/challenges", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	ParticipateChallenge(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]string
	json.NewDecoder(res.Body).Decode(&response)
	assert.Contains(t, response, "challenge_id")
	assert.Equal(t, "processing", response["status"])
}

func TestGetAllChallenges(t *testing.T) {
	challengeID := uuid.New().String()
	challenge := Challenge{
		ID:        challengeID,
		PlayerID:  "player123",
		Amount:    20.01,
		Won:       false,
		Timestamp: time.Now(),
		Status:    "completed",
	}

	challengeMutex.Lock()
	challenges[challengeID] = challenge
	challengeMutex.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/challenges/results", nil)
	w := httptest.NewRecorder()
	GetAllChallenges(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response []Challenge
	json.NewDecoder(res.Body).Decode(&response)
	assert.NotEmpty(t, response)
}

func TestGetChallengeResult(t *testing.T) {
	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/challenges/results/{id}", GetChallengeResult).Methods("GET")

	challengeID := uuid.New().String()
	challenge := Challenge{
		ID:        challengeID,
		PlayerID:  "player123",
		Amount:    20.01,
		Won:       false,
		Timestamp: time.Now(),
		Status:    "completed",
	}

	challengeMutex.Lock()
	challenges[challengeID] = challenge
	challengeMutex.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/challenges/results/"+challengeID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	var result Challenge
	err := json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if result.ID != challengeID {
		t.Fatalf("Expected challenge ID %s, got %s", challengeID, result.ID)
	}
}
