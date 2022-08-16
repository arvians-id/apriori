package mysql

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
)

type userOrderRepository struct {
}

func NewUserOrderRepository() repository.UserOrderRepository {
	return &userOrderRepository{}
}

func (repository *userOrderRepository) FindAll(ctx context.Context, tx *sql.Tx, payloadId string) ([]entity.UserOrder, error) {
	query := "SELECT * FROM user_orders WHERE payload_id = ?"
	queryContext, err := tx.QueryContext(ctx, query, payloadId)
	if err != nil {
		return []entity.UserOrder{}, err
	}

	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var userOrders []entity.UserOrder
	for queryContext.Next() {
		var userOrder entity.UserOrder
		err := queryContext.Scan(&userOrder.IdOrder, &userOrder.PayloadId, &userOrder.Code, &userOrder.Name, &userOrder.Price, &userOrder.Image, &userOrder.Quantity, &userOrder.TotalPriceItem)
		if err != nil {
			return []entity.UserOrder{}, err
		}
		userOrders = append(userOrders, userOrder)
	}

	return userOrders, nil
}
func (repository *userOrderRepository) FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]entity.UserOrderRelationByUserId, error) {
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
	queryContext, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		return []entity.UserOrderRelationByUserId{}, err
	}

	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var userOrders []entity.UserOrderRelationByUserId
	for queryContext.Next() {
		var userOrder entity.UserOrderRelationByUserId
		err := queryContext.Scan(
			&userOrder.IdOrder,
			&userOrder.PayloadId,
			&userOrder.Code,
			&userOrder.Name,
			&userOrder.Price,
			&userOrder.Image,
			&userOrder.Quantity,
			&userOrder.TotalPriceItem,
			&userOrder.OrderId,
			&userOrder.TransactionStatus,
		)
		if err != nil {
			return []entity.UserOrderRelationByUserId{}, err
		}
		userOrders = append(userOrders, userOrder)
	}

	return userOrders, nil
}

func (repository *userOrderRepository) Create(ctx context.Context, tx *sql.Tx, userOrder entity.UserOrder) error {
	query := `INSERT INTO user_orders(payload_id,code,name,price,image,quantity,total_price_item) VALUES(?,?,?,?,?,?,?)`
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
