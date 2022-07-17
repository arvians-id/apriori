package model

type GetPaymentNullableResponse struct {
	IdPayload         int     `json:"id_payload"`
	UserId            *string `json:"user_id"`
	OrderId           *string `json:"order_id"`
	TransactionTime   *string `json:"transaction_time"`
	TransactionStatus *string `json:"transaction_status"`
	TransactionId     *string `json:"transaction_id"`
	StatusCode        *string `json:"status_code"`
	SignatureKey      *string `json:"signature_key"`
	SettlementTime    *string `json:"settlement_time"`
	PaymentType       *string `json:"payment_type"`
	MerchantId        *string `json:"merchant_id"`
	GrossAmount       *string `json:"gross_amount"`
	FraudStatus       *string `json:"fraud_status"`
	BankType          *string `json:"bank_type,omitempty"`
	VANumber          *string `json:"va_number,omitempty"`
	BillerCode        *string `json:"biller_code,omitempty"`
	BillKey           *string `json:"bill_key,omitempty"`
}

type GetPaymentRelationResponse struct {
	IdPayload         int     `json:"id_payload"`
	UserId            *string `json:"user_id"`
	OrderId           *string `json:"order_id"`
	TransactionTime   *string `json:"transaction_time"`
	TransactionStatus *string `json:"transaction_status"`
	TransactionId     *string `json:"transaction_id"`
	StatusCode        *string `json:"status_code"`
	SignatureKey      *string `json:"signature_key"`
	SettlementTime    *string `json:"settlement_time"`
	PaymentType       *string `json:"payment_type"`
	MerchantId        *string `json:"merchant_id"`
	GrossAmount       *string `json:"gross_amount"`
	FraudStatus       *string `json:"fraud_status"`
	BankType          *string `json:"bank_type,omitempty"`
	VANumber          *string `json:"va_number,omitempty"`
	BillerCode        *string `json:"biller_code,omitempty"`
	BillKey           *string `json:"bill_key,omitempty"`
	UserName          *string `json:"user_name,omitempty"`
}
