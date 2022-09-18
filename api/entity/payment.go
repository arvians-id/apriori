package entity

type Payment struct {
	IdPayload         int          `json:"id_payload"`
	UserId            string       `json:"user_id"`
	OrderId           *string      `json:"order_id"`
	TransactionTime   *string      `json:"transaction_time"`
	TransactionStatus *string      `json:"transaction_status"`
	TransactionId     *string      `json:"transaction_id"`
	StatusCode        *string      `json:"status_code"`
	SignatureKey      *string      `json:"signature_key"`
	SettlementTime    *string      `json:"settlement_time"`
	PaymentType       *string      `json:"payment_type"`
	MerchantId        *string      `json:"merchant_id"`
	GrossAmount       *string      `json:"gross_amount"`
	FraudStatus       *string      `json:"fraud_status"`
	BankType          *string      `json:"bank_type"`
	VANumber          *string      `json:"va_number"`
	BillerCode        *string      `json:"biller_code"`
	BillKey           *string      `json:"bill_key"`
	ReceiptNumber     *string      `json:"receipt_number"`
	Address           *string      `json:"address"`
	Courier           *string      `json:"courier"`
	CourierService    *string      `json:"courier_service"`
	User              *User        `json:"user"`
	UserOrder         []*UserOrder `json:"user_order"`
}
