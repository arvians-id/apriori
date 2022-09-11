package model

type GetPaymentResponse struct {
	IdPayload         int    `json:"id_payload"`
	UserId            string `json:"user_id"`
	OrderId           string `json:"order_id"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionId     string `json:"transaction_id"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	SettlementTime    string `json:"settlement_time"`
	PaymentType       string `json:"payment_type"`
	MerchantId        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	BankType          string `json:"bank_type,omitempty"`
	VANumber          string `json:"va_number,omitempty"`
	BillerCode        string `json:"biller_code,omitempty"`
	BillKey           string `json:"bill_key,omitempty"`
	ReceiptNumber     string `json:"receipt_number,omitempty"`
	Address           string `json:"address,omitempty"`
	Courier           string `json:"courier,omitempty"`
	CourierService    string `json:"courier_service,omitempty"`
}

type GetPaymentRelationResponse struct {
	IdPayload         int    `json:"id_payload"`
	UserId            string `json:"user_id"`
	OrderId           string `json:"order_id"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionId     string `json:"transaction_id"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	SettlementTime    string `json:"settlement_time"`
	PaymentType       string `json:"payment_type"`
	MerchantId        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	BankType          string `json:"bank_type,omitempty"`
	VANumber          string `json:"va_number,omitempty"`
	BillerCode        string `json:"biller_code,omitempty"`
	BillKey           string `json:"bill_key,omitempty"`
	ReceiptNumber     string `json:"receipt_number,omitempty"`
	Address           string `json:"address,omitempty"`
	Courier           string `json:"courier,omitempty"`
	CourierService    string `json:"courier_service,omitempty"`
	UserName          string `json:"user_name,omitempty"`
}

type GetPaymentTokenRequest struct {
	GrossAmount    int64    `form:"gross_amount" binding:"required"`
	Items          []string `form:"items" binding:"required"`
	UserId         int      `form:"user_id"`
	CustomerName   string   `form:"customer_name" binding:"required"`
	Address        string   `form:"address" binding:"required"`
	Courier        string   `form:"courier" binding:"required"`
	CourierService string   `form:"courier_service" binding:"required"`
	ShippingCost   int64    `form:"shipping_cost" binding:"required"`
}

type AddReceiptNumberRequest struct {
	OrderId       string `json:"order_id"`
	ReceiptNumber string `json:"receipt_number"`
}
