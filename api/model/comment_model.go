package model

type GetCommentResponse struct {
	IdComment   int    `json:"id_comment"`
	UserOrderId int    `json:"user_order_id"`
	ProductCode string `json:"product_code,omitempty"`
	Description string `json:"description,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Rating      int    `json:"rating,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UserId      int    `json:"user_id,omitempty"`
	UserName    string `json:"user_name,omitempty"`
}

type CreateCommentRequest struct {
	UserOrderId int    `json:"user_order_id"`
	ProductCode string `json:"product_code,omitempty"`
	Description string `json:"description,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Rating      int    `json:"rating" binding:"required"`
	CreatedAt   string `json:"created_at,omitempty"`
}

type GetRatingResponse struct {
	Rating        int `json:"rating"`
	ResultRating  int `json:"result_rating"`
	ResultComment int `json:"result_comment"`
}
