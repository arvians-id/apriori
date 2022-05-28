package model

type CreateAprioriRequest struct {
	Item       string `json:"item"`
	Discount   int32  `json:"discount"`
	Support    int32  `json:"support"`
	Confidence int32  `json:"confidence"`
	RangeDate  string `json:"range_date"`
	Counter    uint64 `json:"counter"`
	CreatedAt  string `json:"created_at"`
}

type GenerateAprioriRequest struct {
	MinimumSupport    int32  `json:"minimum_support" binding:"required,min=10,max=100"`
	MinimumConfidence int32  `json:"minimum_confidence" binding:"required,max=100"`
	MinimumDiscount   int32  `json:"minimum_discount" binding:"required"`
	MaximumDiscount   int32  `json:"maximum_discount" binding:"required,gtefield=MinimumDiscount"`
	StartDate         string `json:"start_date" binding:"required"`
	EndDate           string `json:"end_date" binding:"required"`
}

type GetAprioriResponse struct {
	IdApriori  uint64 `json:"id_apriori"`
	Item       string `json:"item"`
	Discount   int32  `json:"discount"`
	Support    int32  `json:"support"`
	Confidence int32  `json:"confidence"`
	RangeDate  string `json:"range_date"`
	Counter    uint64 `json:"counter"`
	CreatedAt  string `json:"created_at"`
}
