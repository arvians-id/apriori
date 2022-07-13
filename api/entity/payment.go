package entity

type Payment struct {
	IdPayload         int
	UserId            string
	OrderId           string
	TransactionTime   string
	TransactionStatus string
	TransactionId     string
	StatusCode        string
	SignatureKey      string
	SettlementTime    string
	PaymentType       string
	MerchantId        string
	GrossAmount       string
	FraudStatus       string
	BankType          string
	VANumber          string
	BillerCode        string
	BillKey           string
}
type PaymentNullable struct {
	IdPayload         int
	UserId            *string
	OrderId           *string
	TransactionTime   *string
	TransactionStatus *string
	TransactionId     *string
	StatusCode        *string
	SignatureKey      *string
	SettlementTime    *string
	PaymentType       *string
	MerchantId        *string
	GrossAmount       *string
	FraudStatus       *string
	BankType          *string
	VANumber          *string
	BillerCode        *string
	BillKey           *string
}
