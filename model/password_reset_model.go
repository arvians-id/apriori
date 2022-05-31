package model

type CreatePasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdateResetPasswordUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Token    string `json:"token"`
}

type GetPasswordResetResponse struct {
	Email   string `json:"email"`
	Token   string `json:"token"`
	Expired int32  `json:"expired"`
}
