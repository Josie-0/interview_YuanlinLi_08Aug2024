package handlers

import (
	"PaymentProcessingSystem/models"
	"PaymentProcessingSystem/services"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockPaymentProcessor 实现一个模拟的 PaymentProcessor
type MockPaymentProcessor struct{}

func (p *MockPaymentProcessor) ProcessPayment(details models.PaymentDetails) (models.PaymentResponse, error) {
	return models.PaymentResponse{
		Status:        "success",
		TransactionID: "mock123456789",
	}, nil
}

// MockPaymentService 实现 PaymentService 接口
type MockPaymentService struct {
	processors map[string]services.PaymentProcessor // 修正：添加缺失的 "]"
}

func (s *MockPaymentService) ProcessPayment(request *models.PaymentRequest) (models.PaymentResponse, error) {
	processor, exists := s.processors[request.PaymentMethod]
	if !exists {
		return models.PaymentResponse{
			Status:       "failed",
			ErrorMessage: "payment method not supported",
		}, fmt.Errorf("payment method not supported")
	}

	details := models.PaymentDetails{
		PaymentMethod: request.PaymentMethod,
		Amount:        request.Amount,
	}

	response, err := processor.ProcessPayment(details)
	if err != nil {
		return models.PaymentResponse{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, err
	}

	return response, nil
}

func TestHandlePayment(t *testing.T) {
	mockService := &MockPaymentService{
		processors: map[string]services.PaymentProcessor{
			"credit_card": &MockPaymentProcessor{},
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandlePaymentHandler(mockService, w, r)
	})

	reqBody, _ := json.Marshal(models.PaymentRequest{
		PaymentMethod: "credit_card",
		Amount:        100,
	})
	req := httptest.NewRequest("POST", "/payments", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

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

func TestGetPaymentByID(t *testing.T) {
	handler := http.HandlerFunc(GetPaymentByID)

	req := httptest.NewRequest("GET", "/payments?id=123", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", resp.Status)
	}

	var response models.PaymentDetails
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.TransactionID != "123" {
		t.Fatalf("expected transaction ID 123, got %v", response.TransactionID)
	}
}
