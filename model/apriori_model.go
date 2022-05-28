package model

type CreateAprioriRequest struct {
	Item       string  `json:"item"`
	Discount   float64 `json:"discount"`
	Support    float64 `json:"support"`
	Confidence float64 `json:"confidence"`
	RangeDate  string  `json:"range_date"`
	CreatedAt  string  `json:"created_at"`
}

type GetAprioriResponse struct {
	IdApriori  uint64  `json:"id_apriori"`
	Code       string  `json:"code"`
	Item       string  `json:"item"`
	Discount   float64 `json:"discount"`
	Support    float64 `json:"support"`
	Confidence float64 `json:"confidence"`
	RangeDate  string  `json:"range_date"`
	IsActive   bool    `json:"is_active"`
	CreatedAt  string  `json:"created_at"`
}

type GenerateAprioriRequest struct {
	MinimumSupport    int32  `json:"minimum_support" binding:"required,min=10,max=100"`
	MinimumConfidence int32  `json:"minimum_confidence" binding:"required,max=100"`
	MinimumDiscount   int32  `json:"minimum_discount" binding:"required"`
	MaximumDiscount   int32  `json:"maximum_discount" binding:"required,gtefield=MinimumDiscount"`
	StartDate         string `json:"start_date" binding:"required"`
	EndDate           string `json:"end_date" binding:"required"`
}

type GetGenerateAprioriResponse struct {
	ItemSet     []string `json:"item_set"`
	Support     float64  `json:"support"`
	Iterate     int32    `json:"iterate"`
	Transaction int32    `json:"transaction"`
	Confidence  float64  `json:"confidence,omitempty"`
	Discount    float64  `json:"discount,omitempty"`
	Description string   `json:"description,omitempty"`
	RangeDate   string   `json:"range_date"`
}
