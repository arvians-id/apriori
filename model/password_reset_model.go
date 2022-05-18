package model

type CreatePasswordResetRequest struct {
	Email string `json:"email"`
}

type UpdateResetPasswordUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type GetPasswordResetResponse struct {
	Email   string `json:"email"`
	Token   string `json:"token"`
	Expired int32  `json:"expired"`
}
