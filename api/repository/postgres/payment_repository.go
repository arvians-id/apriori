package postgres

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
	"log"
)

type PaymentRepositoryImpl struct {
}

func NewPaymentRepository() repository.PaymentRepository {
	return &PaymentRepositoryImpl{}
}

func (repository *PaymentRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.PaymentRelation, error) {
	query := `SELECT payloads.*,users.name 
			  FROM payloads
			  	LEFT JOIN users ON users.id_user = payloads.user_id
			  ORDER BY payloads.settlement_time DESC, payloads.bank_type DESC`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(rows)

	var payments []*entity.PaymentRelation
	for rows.Next() {
		var payment entity.PaymentRelation
		err := rows.Scan(
			&payment.Payment.IdPayload,
			&payment.Payment.UserId,
			&payment.Payment.OrderId,
			&payment.Payment.TransactionTime,
			&payment.Payment.TransactionStatus,
			&payment.Payment.TransactionId,
			&payment.Payment.StatusCode,
			&payment.Payment.SignatureKey,
			&payment.Payment.SettlementTime,
			&payment.Payment.PaymentType,
			&payment.Payment.MerchantId,
			&payment.Payment.GrossAmount,
			&payment.Payment.FraudStatus,
			&payment.Payment.BankType,
			&payment.Payment.VANumber,
			&payment.Payment.BillerCode,
			&payment.Payment.BillKey,
			&payment.Payment.ReceiptNumber,
			&payment.Payment.Address,
			&payment.Payment.Courier,
			&payment.Payment.CourierService,
			&payment.UserName,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, &payment)
	}

	return payments, nil
}

func (repository *PaymentRepositoryImpl) FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*entity.Payment, error) {
	query := `SELECT * FROM payloads 
			  WHERE user_id = $1 
			  ORDER BY settlement_time DESC, bank_type DESC`
	rows, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(rows)

	var payments []*entity.Payment
	for rows.Next() {
		var payment entity.Payment
		err := rows.Scan(
			&payment.IdPayload,
			&payment.UserId,
			&payment.OrderId,
			&payment.TransactionTime,
			&payment.TransactionStatus,
			&payment.TransactionId,
			&payment.StatusCode,
			&payment.SignatureKey,
			&payment.SettlementTime,
			&payment.PaymentType,
			&payment.MerchantId,
			&payment.GrossAmount,
			&payment.FraudStatus,
			&payment.BankType,
			&payment.VANumber,
			&payment.BillerCode,
			&payment.BillKey,
			&payment.ReceiptNumber,
			&payment.Address,
			&payment.Courier,
			&payment.CourierService,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, &payment)
	}

	return payments, nil
}

func (repository *PaymentRepositoryImpl) FindByOrderId(ctx context.Context, tx *sql.Tx, orderId string) (*entity.Payment, error) {
	query := "SELECT * FROM payloads WHERE order_id = $1"
	row := tx.QueryRowContext(ctx, query, orderId)

	var payment entity.Payment
	err := row.Scan(
		&payment.IdPayload,
		&payment.UserId,
		&payment.OrderId,
		&payment.TransactionTime,
		&payment.TransactionStatus,
		&payment.TransactionId,
		&payment.StatusCode,
		&payment.SignatureKey,
		&payment.SettlementTime,
		&payment.PaymentType,
		&payment.MerchantId,
		&payment.GrossAmount,
		&payment.FraudStatus,
		&payment.BankType,
		&payment.VANumber,
		&payment.BillerCode,
		&payment.BillKey,
		&payment.ReceiptNumber,
		&payment.Address,
		&payment.Courier,
		&payment.CourierService,
	)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (repository *PaymentRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, payment *entity.Payment) (*entity.Payment, error) {
	id := 0
	query := `INSERT INTO payloads(user_id,order_id,transaction_time,transaction_status,transaction_id,status_code,signature_key,settlement_time,payment_type,merchant_id,gross_amount,fraud_status,bank_type,va_number,biller_code,bill_key,receipt_number,address,courier,courier_service) 
			  VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)  RETURNING id_payload`
	row := tx.QueryRowContext(
		ctx,
		query,
		payment.UserId,
		payment.OrderId,
		payment.TransactionTime,
		payment.TransactionStatus,
		payment.TransactionId,
		payment.StatusCode,
		payment.SignatureKey,
		payment.SettlementTime,
		payment.PaymentType,
		payment.MerchantId,
		payment.GrossAmount,
		payment.FraudStatus,
		payment.BankType,
		payment.VANumber,
		payment.BillerCode,
		payment.BillKey,
		payment.ReceiptNumber,
		payment.Address,
		payment.Courier,
		payment.CourierService,
	)
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	payment.IdPayload = id

	return payment, nil
}

func (repository *PaymentRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, payment *entity.Payment) error {
	query := `UPDATE payloads 
         	  SET user_id = $1,
         	      order_id = $2,
         	      transaction_time = $3,
         	      transaction_status = $4,
         	      transaction_id = $5,
				  status_code = $6,
				  signature_key = $7,
				  settlement_time = $8,
         	      payment_type = $9,
         	      merchant_id = $10,
         	      gross_amount = $11,
         	      fraud_status = $12,
		 	      bank_type = $13,
         	      va_number = $14,
         	      biller_code = $15,
         	      bill_key = $16
			  WHERE order_id = $17`
	_, err := tx.ExecContext(
		ctx,
		query,
		payment.UserId,
		payment.OrderId,
		payment.TransactionTime,
		payment.TransactionStatus,
		payment.TransactionId,
		payment.StatusCode,
		payment.SignatureKey,
		payment.SettlementTime,
		payment.PaymentType,
		payment.MerchantId,
		payment.GrossAmount,
		payment.FraudStatus,
		payment.BankType,
		payment.VANumber,
		payment.BillerCode,
		payment.BillKey,
		payment.OrderId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *PaymentRepositoryImpl) UpdateReceiptNumber(ctx context.Context, tx *sql.Tx, payment *entity.Payment) error {
	query := `UPDATE payloads SET receipt_number = $1 WHERE order_id = $2`
	_, err := tx.ExecContext(ctx, query, payment.ReceiptNumber, payment.OrderId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *PaymentRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, orderId string) error {
	query := "DELETE FROM payloads WHERE order_id = $1"
	_, err := tx.ExecContext(ctx, query, orderId)
	if err != nil {
		return err
	}

	return nil
}
