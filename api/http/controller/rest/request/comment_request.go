package request

type CreateCommentRequest struct {
	UserOrderId int    `json:"user_order_id"`
	ProductCode string `json:"product_code"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	Rating      int    `json:"rating" binding:"required"`
	CreatedAt   string `json:"created_at"`
}
