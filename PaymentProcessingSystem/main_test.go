package main

import (
	"PaymentProcessingSystem/services"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"PaymentProcessingSystem/handlers"
	"PaymentProcessingSystem/models"
	"github.com/gorilla/mux"
)

func setup() {
	// 设置环境变量
	os.Setenv("API_KEY", "test-api-key")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestAPIKeyMiddleware(t *testing.T) {
	router := mux.NewRouter()
	router.Use(APIKeyMiddleware)
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-KEY", "test-api-key")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", rec.Result().Status)
	}
}

func TestHandlePaymentIntegration(t *testing.T) {
	router := mux.NewRouter()
	paymentService := services.NewPaymentService()
	router.HandleFunc("/payments", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlePaymentHandler(paymentService, w, r)
	}).Methods("POST")
	router.Use(APIKeyMiddleware)

	reqBody, _ := json.Marshal(models.PaymentRequest{
		PaymentMethod: "credit_card",
		Amount:        100,
	})
	req := httptest.NewRequest("POST", "/payments", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test-api-key")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", resp.Status)
	}

	var response models.PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "success" {
		t.Fatalf("expected success status, got %v", response.Status)
	}
	if response.TransactionID == "" {
		t.Fatal("expected non-empty transaction ID")
	}
}

func TestGetPaymentByIDIntegration(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/payments/{id:[0-9]+}", handlers.GetPaymentByID).Methods("GET")
	router.Use(APIKeyMiddleware)

	req := httptest.NewRequest("GET", "/payments/123", nil)
	req.Header.Set("X-API-KEY", "test-api-key")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", resp.Status)
	}

	var response models.PaymentDetails
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}
