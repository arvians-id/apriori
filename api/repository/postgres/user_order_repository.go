package postgres

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
	"log"
)

type UserOrderRepositoryImpl struct {
}

func NewUserOrderRepository() repository.UserOrderRepository {
	return &UserOrderRepositoryImpl{}
}

func (repository *UserOrderRepositoryImpl) FindAllByPayloadId(ctx context.Context, tx *sql.Tx, payloadId string) ([]*entity.UserOrder, error) {
	query := "SELECT * FROM user_orders WHERE payload_id = $1"
	rows, err := tx.QueryContext(ctx, query, payloadId)
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

	var userOrders []*entity.UserOrder
	for rows.Next() {
		var userOrder entity.UserOrder
		err := rows.Scan(
			&userOrder.IdOrder,
			&userOrder.PayloadId,
			&userOrder.Code,
			&userOrder.Name,
			&userOrder.Price,
			&userOrder.Image,
			&userOrder.Quantity,
			&userOrder.TotalPriceItem,
		)
		if err != nil {
			return nil, err
		}

		userOrders = append(userOrders, &userOrder)
	}

	return userOrders, nil
}

func (repository *UserOrderRepositoryImpl) FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*entity.UserOrderRelationByUserId, error) {
	query := `SELECT 
				id_order,
			    payload_id,
			    code,
      			name,
			    price,
			    image,
			    quantity,
			    total_price_item,
			    order_id,
			    transaction_status 
			  FROM user_orders uo 
			  	LEFT JOIN payloads p ON p.id_payload = uo.payload_id 
		   	  WHERE p.user_id = $1 AND p.transaction_status = 'settlement'
			  ORDER BY uo.id_order DESC`
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

	var userOrders []*entity.UserOrderRelationByUserId
	for rows.Next() {
		var userOrder entity.UserOrderRelationByUserId
		err := rows.Scan(
			&userOrder.UserOrder.IdOrder,
			&userOrder.UserOrder.PayloadId,
			&userOrder.UserOrder.Code,
			&userOrder.UserOrder.Name,
			&userOrder.UserOrder.Price,
			&userOrder.UserOrder.Image,
			&userOrder.UserOrder.Quantity,
			&userOrder.UserOrder.TotalPriceItem,
			&userOrder.OrderId,
			&userOrder.TransactionStatus,
		)
		if err != nil {
			return nil, err
		}

		userOrders = append(userOrders, &userOrder)
	}

	return userOrders, nil
}

func (repository *UserOrderRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (*entity.UserOrder, error) {
	query := `SELECT * FROM user_orders WHERE id_order = $1`
	row := tx.QueryRowContext(ctx, query, id)

	var userOrder entity.UserOrder
	err := row.Scan(
		&userOrder.IdOrder,
		&userOrder.PayloadId,
		&userOrder.Code,
		&userOrder.Name,
		&userOrder.Price,
		&userOrder.Image,
		&userOrder.Quantity,
		&userOrder.TotalPriceItem,
	)
	if err != nil {
		return nil, err
	}

	return &userOrder, nil
}

func (repository *UserOrderRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, userOrder *entity.UserOrder) error {
	query := `INSERT INTO user_orders(payload_id,code,name,price,image,quantity,total_price_item) VALUES($1,$2,$3,$4,$5,$6,$7)`
	_, err := tx.ExecContext(
		ctx,
		query,
		userOrder.PayloadId,
		userOrder.Code,
		userOrder.Name,
		userOrder.Price,
		userOrder.Image,
		userOrder.Quantity,
		userOrder.TotalPriceItem,
	)
	if err != nil {
		return err
	}

	return nil
}
