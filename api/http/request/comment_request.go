package request

type CreateCommentRequest struct {
	UserOrderId int    `json:"user_order_id"`
	ProductCode string `json:"product_code,omitempty"`
	Description string `json:"description,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Rating      int    `json:"rating" binding:"required"`
	CreatedAt   string `json:"created_at,omitempty"`
}
