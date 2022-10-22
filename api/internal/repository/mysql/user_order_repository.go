package mysql

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"log"
)

type UserOrderRepositoryImpl struct {
}

func NewUserOrderRepository() repository.UserOrderRepository {
	return &UserOrderRepositoryImpl{}
}

func (repository *UserOrderRepositoryImpl) FindAllByPayloadId(ctx context.Context, tx *sql.Tx, payloadId string) ([]*model.UserOrder, error) {
	query := "SELECT * FROM user_orders WHERE payload_id = ?"
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

	var userOrders []*model.UserOrder
	for rows.Next() {
		var userOrder model.UserOrder
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

func (repository *UserOrderRepositoryImpl) FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*model.UserOrder, error) {
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
		   	  WHERE p.user_id = ? AND p.transaction_status = 'settlement'
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

	var userOrders []*model.UserOrder
	for rows.Next() {
		userOrder := model.UserOrder{
			Payment: &model.Payment{},
		}
		err := rows.Scan(
			&userOrder.IdOrder,
			&userOrder.PayloadId,
			&userOrder.Code,
			&userOrder.Name,
			&userOrder.Price,
			&userOrder.Image,
			&userOrder.Quantity,
			&userOrder.TotalPriceItem,
			&userOrder.Payment.OrderId,
			&userOrder.Payment.TransactionStatus,
		)
		if err != nil {
			return nil, err
		}

		userOrders = append(userOrders, &userOrder)
	}

	return userOrders, nil
}

func (repository *UserOrderRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (*model.UserOrder, error) {
	query := `SELECT * FROM user_orders WHERE id_order = ?`
	row := tx.QueryRowContext(ctx, query, id)

	var userOrder model.UserOrder
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
		return &model.UserOrder{}, err
	}

	return &userOrder, nil
}

func (repository *UserOrderRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, userOrder *model.UserOrder) (*model.UserOrder, error) {
	query := `INSERT INTO user_orders(payload_id,code,name,price,image,quantity,total_price_item) VALUES(?,?,?,?,?,?,?)`
	row, err := tx.ExecContext(
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
		return nil, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}

	userOrder.IdOrder = int(id)

	return userOrder, nil
}
