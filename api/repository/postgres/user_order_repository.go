package postgres

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
	query := "SELECT * FROM user_orders WHERE payload_id = $1"
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

func (repository *userOrderRepository) Create(ctx context.Context, tx *sql.Tx, userOrder entity.UserOrder) error {
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
