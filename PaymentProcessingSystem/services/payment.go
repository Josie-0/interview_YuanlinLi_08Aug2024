package services

import (
	"PaymentProcessingSystem/models"
)

// PaymentProcessor 定义支付处理接口
type PaymentProcessor interface {
	ProcessPayment(details models.PaymentDetails) (models.PaymentResponse, error)
}

// 其他支付处理器实现
type CreditCardProcessor struct{}

func (p *CreditCardProcessor) ProcessPayment(details models.PaymentDetails) (models.PaymentResponse, error) {
	// 模拟调用信用卡支付网关
	return models.PaymentResponse{
		Status:        "success",
		TransactionID: "cc123456789",
	}, nil
}

type BankTransferProcessor struct{}

func (p *BankTransferProcessor) ProcessPayment(details models.PaymentDetails) (models.PaymentResponse, error) {
	// 模拟调用银行转账服务
	return models.PaymentResponse{
		Status:        "success",
		TransactionID: "bt123456789",
	}, nil
}

type ThirdPartyProcessor struct{}

func (p *ThirdPartyProcessor) ProcessPayment(details models.PaymentDetails) (models.PaymentResponse, error) {
	// 模拟调用第三方支付平台
	return models.PaymentResponse{
		Status:        "success",
		TransactionID: "tp123456789",
	}, nil
}

type BlockchainProcessor struct{}

func (p *BlockchainProcessor) ProcessPayment(details models.PaymentDetails) (models.PaymentResponse, error) {
	// 模拟调用区块链支付网关
	return models.PaymentResponse{
		Status:        "success",
		TransactionID: "bc123456789",
	}, nil
}
