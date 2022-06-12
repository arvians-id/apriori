package entity

type PasswordReset struct {
	Email   string
	Token   string
	Expired int32
}
