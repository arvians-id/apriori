package model

type GetCommentResponse struct {
	IdComment   int    `json:"id_comment,omitempty"`
	UserOrderId int    `json:"user_order_id,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	Description string `json:"description,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Rating      int    `json:"rating,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UserId      int    `json:"user_id,omitempty"`
	UserName    string `json:"user_name"`
}

type CreateCommentRequest struct {
	UserOrderId int    `json:"user_order_id,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	Description string `json:"description,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Rating      int    `json:"rating,omitempty" binding:"required"`
	CreatedAt   string `json:"created_at,omitempty"`
	UserId      int    `json:"user_id,omitempty"`
	UserName    string `json:"user_name"`
}
