package main

import (
	"log"
	"net/http"
	"os"

	"PaymentProcessingSystem/handlers"
	"PaymentProcessingSystem/services"
	"github.com/gorilla/mux"
)

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		expectedAPIKey := os.Getenv("API_KEY")

		if expectedAPIKey == "" {
			log.Fatal("API_KEY environment variable not set")
		}

		if apiKey != expectedAPIKey {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()

	paymentService := services.NewPaymentService()

	router.HandleFunc("/payments", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlePaymentHandler(paymentService, w, r)
	}).Methods("POST")

	router.HandleFunc("/payments/{id:[a-zA-Z0-9]+}", handlers.GetPaymentByID).Methods("GET")

	router.Use(APIKeyMiddleware)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
