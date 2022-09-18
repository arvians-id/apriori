package entity

type PasswordReset struct {
	Email   string `json:"email"`
	Token   string `json:"token"`
	Expired int32  `json:"expired"`
}
