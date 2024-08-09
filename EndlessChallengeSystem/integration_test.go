// integration_test.go
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/challenges/results", GetAllChallenges).Methods("GET")
	router.HandleFunc("/challenges/results/{id}", GetChallengeResult).Methods("GET")
	router.HandleFunc("/challenges", ParticipateChallenge).Methods("POST")

	// Test Participate Challenge
	reqBody := `{"player_id": "player123", "amount": 20.01}`
	req := httptest.NewRequest(http.MethodPost, "/challenges", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]string
	json.NewDecoder(res.Body).Decode(&response)
	challengeID, ok := response["challenge_id"]
	assert.True(t, ok)

	// Wait for challenge processing
	time.Sleep(35 * time.Second)

	// Test Get Challenge Result
	req = httptest.NewRequest(http.MethodGet, "/challenges/results/"+challengeID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res = w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var result Challenge
	json.NewDecoder(res.Body).Decode(&result)
	assert.Equal(t, challengeID, result.ID)
}
