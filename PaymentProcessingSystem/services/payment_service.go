package services

import (
	"PaymentProcessingSystem/models"
	"fmt"
)

// PaymentService 负责选择合适的支付处理器
type PaymentService interface {
	ProcessPayment(request *models.PaymentRequest) (models.PaymentResponse, error)
}

// ConcretePaymentService 实现 PaymentService 接口
type ConcretePaymentService struct {
	processors map[string]PaymentProcessor
}

// NewPaymentService 创建新的支付服务实例
func NewPaymentService() *ConcretePaymentService {
	return &ConcretePaymentService{
		processors: map[string]PaymentProcessor{
			"credit_card":   &CreditCardProcessor{},
			"bank_transfer": &BankTransferProcessor{},
			"third_party":   &ThirdPartyProcessor{},
			"blockchain":    &BlockchainProcessor{},
		},
	}
}

// ProcessPayment 处理支付请求
func (s *ConcretePaymentService) ProcessPayment(request *models.PaymentRequest) (models.PaymentResponse, error) {
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
