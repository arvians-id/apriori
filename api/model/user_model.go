package model

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,max=20"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Address  string `json:"address" binding:"required,max=100"`
	Phone    string `json:"phone" binding:"required,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	IdUser   uint64 `json:"id_user"`
	Role     int    `json:"role" binding:"omitempty,min=1,max=2"`
	Name     string `json:"name" binding:"required,max=20"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Address  string `json:"address" binding:"required,max=100"`
	Phone    string `json:"phone" binding:"required,max=20"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

type GetUserCredentialRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GetUserResponse struct {
	IdUser    uint64 `json:"id_user"`
	Role      int    `json:"role"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
