package entity

type PasswordReset struct {
	Email   string `json:"email"`
	Token   string `json:"token"`
	Expired int64  `json:"expired"`
}
