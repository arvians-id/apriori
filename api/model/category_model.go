package model

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,max=20"`
}

type UpdateCategoryRequest struct {
	IdCategory int    `json:"id_category"`
	Name       string `json:"name" binding:"required,max=20"`
}

type GetCategoryResponse struct {
	IdCategory int    `json:"id_category"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
