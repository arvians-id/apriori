package mysql

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
	"errors"
)

type paymentRepository struct {
}

func NewPaymentRepository() repository.PaymentRepository {
	return &paymentRepository{}
}

func (repository *paymentRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.PaymentRelation, error) {
	query := `SELECT payloads.*,users.name FROM payloads
			  LEFT JOIN users ON users.id_user = payloads.user_id
			  ORDER BY payloads.settlement_time DESC, payloads.bank_type DESC`
	queryContext, err := tx.QueryContext(ctx, query)
	if err != nil {
		return []entity.PaymentRelation{}, err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var payments []entity.PaymentRelation
	for queryContext.Next() {
		var payment entity.PaymentRelation
		err := queryContext.Scan(
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
			&payment.UserName,
		)
		if err != nil {
			return []entity.PaymentRelation{}, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (repository *paymentRepository) FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]entity.PaymentNullable, error) {
	query := "SELECT * FROM payloads WHERE user_id = ? ORDER BY settlement_time DESC, bank_type DESC"
	queryContext, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		return []entity.PaymentNullable{}, err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var payments []entity.PaymentNullable
	for queryContext.Next() {
		var payment entity.PaymentNullable
		err := queryContext.Scan(
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
		)
		if err != nil {
			return []entity.PaymentNullable{}, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (repository *paymentRepository) Create(ctx context.Context, tx *sql.Tx, payment entity.Payment) (entity.Payment, error) {
	query := `INSERT INTO payloads(user_id,order_id,transaction_time,transaction_status,transaction_id,status_code,signature_key,settlement_time,payment_type,merchant_id,gross_amount,fraud_status,bank_type,va_number,biller_code,bill_key) 
			  VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	row, err := tx.ExecContext(
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
		payment.MerchantId, payment.GrossAmount,
		payment.FraudStatus,
		payment.BankType,
		payment.VANumber,
		payment.BillerCode,
		payment.BillKey,
	)
	if err != nil {
		return entity.Payment{}, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return entity.Payment{}, err
	}

	payment.IdPayload = int(id)

	return payment, nil
}

func (repository *paymentRepository) FindByOrderId(ctx context.Context, tx *sql.Tx, orderId string) (entity.PaymentNullable, error) {
	query := "SELECT * FROM payloads WHERE order_id = ?"
	queryContext, err := tx.QueryContext(ctx, query, orderId)
	if err != nil {
		return entity.PaymentNullable{}, err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var payment entity.PaymentNullable
	if queryContext.Next() {
		err := queryContext.Scan(
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
		)
		if err != nil {
			return entity.PaymentNullable{}, err
		}

		return payment, nil
	}

	return payment, errors.New("payment not found")
}

func (repository *paymentRepository) Update(ctx context.Context, tx *sql.Tx, payment entity.Payment) error {
	query := `UPDATE payloads 
         	  SET 
		 	      user_id = ?,
         	      order_id = ?,
         	      transaction_time = ?,
         	      transaction_status = ?,
         	      transaction_id = ?,
				  status_code = ?,
				  signature_key = ?,
				  settlement_time = ?,
         	      payment_type = ?,
         	      merchant_id = ?,
         	      gross_amount = ?,
         	      fraud_status = ?,
		 	      bank_type = ?,
         	      va_number = ?,
         	      biller_code = ?,
         	      bill_key = ?
			  WHERE order_id = ?`
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

func (repository *paymentRepository) AddReceiptNumber(ctx context.Context, tx *sql.Tx, payment entity.Payment) error {
	query := `UPDATE payloads SET receipt_number = ? WHERE order_id = ?`
	_, err := tx.ExecContext(ctx, query, payment.ReceiptNumber, payment.OrderId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *paymentRepository) Delete(ctx context.Context, tx *sql.Tx, orderId string) error {
	query := "DELETE FROM payloads WHERE order_id = ?"
	_, err := tx.ExecContext(ctx, query, orderId)
	if err != nil {
		return err
	}

	return nil
}
