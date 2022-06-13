package model

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,max=20"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	IdUser   uint64 `json:"id_user"`
	Name     string `json:"name" binding:"required,max=20"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

type GetUserCredentialRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GetUserResponse struct {
	IdUser    uint64 `json:"id_user"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
