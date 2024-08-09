package models

type PaymentRequest struct {
	PaymentMethod string                 `json:"payment_method"` // 支付方式：信用卡、银行转账、第三方支付、区块链支付
	Amount        float64                `json:"amount"`         // 支付金额，以最小单位（例如分）表示
	Details       map[string]interface{} `json:"details"`        // 支付详细信息，可以根据具体支付方式存储不同的数据
}

type PaymentResponse struct {
	Status        string `json:"status"`                  // 支付状态：成功或失败
	TransactionID string `json:"transaction_id"`          // 交易 ID
	ErrorMessage  string `json:"error_message,omitempty"` // 错误信息（如果有）
}

type PaymentDetails struct {
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	TransactionID string  `json:"transaction_id"`
}
