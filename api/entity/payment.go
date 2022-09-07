package entity

import (
	"apriori/model"
	"database/sql"
)

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

func (payment *Payment) ToPaymentResponse() *model.GetPaymentResponse {
	return &model.GetPaymentResponse{
		IdPayload:         payment.IdPayload,
		UserId:            payment.UserId.String,
		OrderId:           payment.OrderId.String,
		TransactionTime:   payment.TransactionTime.String,
		TransactionStatus: payment.TransactionStatus.String,
		TransactionId:     payment.TransactionId.String,
		StatusCode:        payment.StatusCode.String,
		SignatureKey:      payment.SignatureKey.String,
		SettlementTime:    payment.SettlementTime.String,
		PaymentType:       payment.PaymentType.String,
		MerchantId:        payment.MerchantId.String,
		GrossAmount:       payment.GrossAmount.String,
		FraudStatus:       payment.FraudStatus.String,
		BankType:          payment.BankType.String,
		VANumber:          payment.VANumber.String,
		BillerCode:        payment.BillerCode.String,
		BillKey:           payment.BillKey.String,
		ReceiptNumber:     payment.ReceiptNumber.String,
		Address:           payment.Address.String,
		Courier:           payment.Courier.String,
		CourierService:    payment.CourierService.String,
	}
}

type PaymentRelation struct {
	Payment  Payment
	UserName sql.NullString
}

func (payment *PaymentRelation) ToPaymentRelationResponse() *model.GetPaymentRelationResponse {
	return &model.GetPaymentRelationResponse{
		IdPayload:         payment.Payment.IdPayload,
		UserId:            payment.Payment.UserId.String,
		OrderId:           payment.Payment.OrderId.String,
		TransactionTime:   payment.Payment.TransactionTime.String,
		TransactionStatus: payment.Payment.TransactionStatus.String,
		TransactionId:     payment.Payment.TransactionId.String,
		StatusCode:        payment.Payment.StatusCode.String,
		SignatureKey:      payment.Payment.SignatureKey.String,
		SettlementTime:    payment.Payment.SettlementTime.String,
		PaymentType:       payment.Payment.PaymentType.String,
		MerchantId:        payment.Payment.MerchantId.String,
		GrossAmount:       payment.Payment.GrossAmount.String,
		FraudStatus:       payment.Payment.FraudStatus.String,
		BankType:          payment.Payment.BankType.String,
		VANumber:          payment.Payment.VANumber.String,
		BillerCode:        payment.Payment.BillerCode.String,
		BillKey:           payment.Payment.BillKey.String,
		ReceiptNumber:     payment.Payment.ReceiptNumber.String,
		Address:           payment.Payment.Address.String,
		Courier:           payment.Payment.Courier.String,
		CourierService:    payment.Payment.CourierService.String,
		UserName:          payment.UserName.String,
	}
}
