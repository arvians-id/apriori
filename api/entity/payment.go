package entity

type Payment struct {
	IdPayload         int
	UserId            int
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
