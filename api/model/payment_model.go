package model

type GetPaymentTokenRequest struct {
	GrossAmount    int64    `form:"gross_amount" binding:"required"`
	Items          []string `form:"items" binding:"required"`
	UserId         int      `form:"user_id"`
	CustomerName   string   `form:"customer_name" binding:"required"`
	Address        string   `form:"address" binding:"required"`
	Courier        string   `form:"courier" binding:"required"`
	CourierService string   `form:"courier_service" binding:"required"`
	ShippingCost   int64    `form:"shipping_cost" binding:"required"`
}

type AddReceiptNumberRequest struct {
	OrderId       string `json:"order_id"`
	ReceiptNumber string `json:"receipt_number"`
}
