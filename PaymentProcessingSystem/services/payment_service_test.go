package services

import (
	"PaymentProcessingSystem/models"
	"testing"
)

func TestProcessPayment_CreditCard(t *testing.T) {
	service := NewPaymentService()
	request := &models.PaymentRequest{
		PaymentMethod: "credit_card",
		Amount:        100,
	}

	response, err := service.ProcessPayment(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if response.Status != "success" {
		t.Fatalf("expected success status, got %v", response.Status)
	}
	if response.TransactionID == "" {
		t.Fatal("expected non-empty transaction ID")
	}
}

func TestProcessPayment_BankTransfer(t *testing.T) {
	service := NewPaymentService()
	request := &models.PaymentRequest{
		PaymentMethod: "bank_transfer",
		Amount:        200,
	}

	response, err := service.ProcessPayment(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if response.Status != "success" {
		t.Fatalf("expected success status, got %v", response.Status)
	}
	if response.TransactionID == "" {
		t.Fatal("expected non-empty transaction ID")
	}
}

func TestProcessPayment_InvalidMethod(t *testing.T) {
	service := NewPaymentService()
	request := &models.PaymentRequest{
		PaymentMethod: "invalid_method",
		Amount:        100,
	}

	response, err := service.ProcessPayment(request)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if response.Status != "failed" {
		t.Fatalf("expected failed status, got %v", response.Status)
	}
	if response.ErrorMessage != "payment method not supported" {
		t.Fatalf("expected 'payment method not supported' error message, got %v", response.ErrorMessage)
	}
}
