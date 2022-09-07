package entity

import "apriori/model"

type UserOrder struct {
	IdOrder        int
	PayloadId      int
	Code           string
	Name           string
	Price          int64
	Image          string
	Quantity       int
	TotalPriceItem int64
}

func (userOrder *UserOrder) ToUserOrderResponse() *model.GetUserOrderResponse {
	return &model.GetUserOrderResponse{
		IdOrder:        userOrder.IdOrder,
		PayloadId:      userOrder.PayloadId,
		Code:           userOrder.Code,
		Name:           userOrder.Name,
		Price:          userOrder.Price,
		Image:          userOrder.Image,
		Quantity:       userOrder.Quantity,
		TotalPriceItem: userOrder.TotalPriceItem,
	}
}

type UserOrderRelationByUserId struct {
	UserOrder         UserOrder
	OrderId           string
	TransactionStatus string
}

func (userOrder *UserOrderRelationByUserId) ToUserOrderRelationByUserIdResponse() *model.GetUserOrderRelationByUserIdResponse {
	return &model.GetUserOrderRelationByUserIdResponse{
		IdOrder:           userOrder.UserOrder.IdOrder,
		PayloadId:         userOrder.UserOrder.PayloadId,
		Code:              userOrder.UserOrder.Code,
		Name:              userOrder.UserOrder.Name,
		Price:             userOrder.UserOrder.Price,
		Image:             userOrder.UserOrder.Image,
		Quantity:          userOrder.UserOrder.Quantity,
		TotalPriceItem:    userOrder.UserOrder.TotalPriceItem,
		OrderId:           userOrder.OrderId,
		TransactionStatus: userOrder.TransactionStatus,
	}
}
