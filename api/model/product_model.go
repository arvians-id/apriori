package model

type CreateProductRequest struct {
	Code        string `form:"code" binding:"required,min=2,max=10"`
	Name        string `form:"name" binding:"required,min=6,max=100"`
	Description string `form:"description" binding:"omitempty,max=100"`
	Price       int    `form:"price" binding:"min=0"`
	Category    string `form:"category" binding:"required,max=100"`
	Mass        int    `form:"mass"`
	Image       string `form:"-"`
}

type UpdateProductRequest struct {
	IdProduct   int    `form:"-"`
	Code        string `form:"-"`
	Name        string `form:"name" binding:"required,min=6,max=100"`
	Description string `form:"description" binding:"omitempty,max=100"`
	Price       int    `form:"price" binding:"min=0"`
	Category    string `form:"category" binding:"required,max=100"`
	IsEmpty     bool   `form:"is_empty"`
	Mass        int    `form:"mass"`
	Image       string `form:"-"`
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
