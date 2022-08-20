package entity

import "database/sql"

type Payment struct {
	IdPayload         int
	UserId            sql.NullString
	OrderId           sql.NullString
	TransactionTime   sql.NullString
	TransactionStatus sql.NullString
	TransactionId     sql.NullString
	StatusCode        sql.NullString
	SignatureKey      sql.NullString
	SettlementTime    sql.NullString
	PaymentType       sql.NullString
	MerchantId        sql.NullString
	GrossAmount       sql.NullString
	FraudStatus       sql.NullString
	BankType          sql.NullString
	VANumber          sql.NullString
	BillerCode        sql.NullString
	BillKey           sql.NullString
	ReceiptNumber     sql.NullString
	Address           sql.NullString
	Courier           sql.NullString
	CourierService    sql.NullString
}

type PaymentRelation struct {
	Payment  Payment
	UserName sql.NullString
}
