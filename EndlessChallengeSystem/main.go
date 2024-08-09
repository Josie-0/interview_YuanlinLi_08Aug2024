package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/challenges/results", GetAllChallenges).Methods("GET")
	router.HandleFunc("/challenges/results/{id}", GetChallengeResult).Methods("GET")
	router.HandleFunc("/challenges", ParticipateChallenge).Methods("POST")

	log.Println("Server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
