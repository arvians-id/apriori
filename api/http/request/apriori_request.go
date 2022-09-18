package request

type CreateAprioriRequest struct {
	Item       string  `json:"item"`
	Discount   float64 `json:"discount"`
	Support    float64 `json:"support"`
	Confidence float64 `json:"confidence"`
	RangeDate  string  `json:"range_date"`
	CreatedAt  string  `json:"created_at"`
}

type UpdateAprioriRequest struct {
	IdApriori   int
	Code        string
	Description string
	Image       string
}

type GenerateAprioriRequest struct {
	MinimumSupport    float64 `json:"minimum_support" binding:"required,max=100"`
	MinimumConfidence float64 `json:"minimum_confidence" binding:"required,max=100"`
	MinimumDiscount   int32   `json:"minimum_discount" binding:"required"`
	MaximumDiscount   int32   `json:"maximum_discount" binding:"required,gtefield=MinimumDiscount"`
	StartDate         string  `json:"start_date" binding:"required"`
	EndDate           string  `json:"end_date" binding:"required"`
}
