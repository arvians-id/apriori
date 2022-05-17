package model

type CreatePasswordResetRequest struct {
	Email string `json:"email"`
}

type GetPasswordResetResponse struct {
	Email   string `json:"email"`
	Token   string `json:"token"`
	Expired int32  `json:"expired"`
}
