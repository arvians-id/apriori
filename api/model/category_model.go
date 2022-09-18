package model

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,max=20"`
}

type UpdateCategoryRequest struct {
	IdCategory int    `json:"id_category"`
	Name       string `json:"name" binding:"required,max=20"`
}
