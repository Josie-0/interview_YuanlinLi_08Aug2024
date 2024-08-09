package handlers

import (
	"PaymentProcessingSystem/models"
	"PaymentProcessingSystem/services"
	"encoding/json"
	"net/http"
)

// HandlePaymentHandler 处理支付请求的处理函数
func HandlePaymentHandler(paymentService services.PaymentService, w http.ResponseWriter, r *http.Request) {
	var request models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	response, err := paymentService.ProcessPayment(&request)
	if err != nil {
		http.Error(w, "Payment processing failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPaymentByID 获取支付详情
func GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	response := models.PaymentDetails{
		PaymentMethod: "credit_card",
		Amount:        100.00,
		TransactionID: id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
