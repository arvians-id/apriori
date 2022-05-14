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
