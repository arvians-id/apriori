package model

type CreateProductRequest struct {
	Code        string `json:"code" binding:"required,min=2,max=10"`
	Name        string `json:"name" binding:"required,min=6,max=100"`
	Description string `json:"description" binding:"omitempty,max=100"`
	Price       int    `json:"price" binding:"min=0"`
	Category    string `json:"category" binding:"required,max=100"`
	Mass        int    `json:"mass"`
	Image       string `json:"image"`
}

type UpdateProductRequest struct {
	IdProduct   int    `json:"id_product"`
	Code        string `json:"code" binding:"required,min=2,max=10"`
	Name        string `json:"name" binding:"required,min=6,max=100"`
	Description string `json:"description" binding:"omitempty,max=100"`
	Price       int    `json:"price" binding:"min=0"`
	Category    string `json:"category" binding:"required,max=100"`
	IsEmpty     bool   `json:"is_empty"`
	Mass        int    `json:"mass"`
	Image       string `json:"image"`
}

type GetProductResponse struct {
	IdProduct   int    `json:"id_product"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Category    string `json:"category"`
	IsEmpty     bool   `json:"is_empty"`
	Mass        int    `json:"mass"`
	Image       string `json:"image"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GetProductRecommendationResponse struct {
	AprioriId          int     `json:"apriori_id"`
	AprioriCode        string  `json:"apriori_code"`
	AprioriItem        string  `json:"apriori_item"`
	AprioriDiscount    float64 `json:"apriori_discount"`
	ProductTotalPrice  int     `json:"product_total_price"`
	PriceAfterDiscount int     `json:"price_discount"`
	Image              string  `json:"apriori_image,omitempty"`
	Mass               int     `json:"mass,omitempty"`
	Description        string  `json:"apriori_description,omitempty"`
}

type GetProductTransactionResponse struct {
	Code        string  `json:"code"`
	ProductName string  `json:"product_name"`
	Transaction int32   `json:"transaction"`
	Support     float32 `json:"support"`
}
