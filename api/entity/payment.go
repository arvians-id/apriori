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
	ReceiptNumber     string
	Address           string
	Courier           string
	CourierService    string
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
	ReceiptNumber     *string
	Address           *string
	Courier           *string
	CourierService    *string
}

type PaymentRelation struct {
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
	ReceiptNumber     *string
	Address           *string
	Courier           *string
	CourierService    *string
	UserName          *string
}
