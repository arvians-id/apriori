package entity

import "apriori/model"

type PasswordReset struct {
	Email   string
	Token   string
	Expired int32
}

func (reset *PasswordReset) ToPasswordResetResponse() *model.GetPasswordResetResponse {
	return &model.GetPasswordResetResponse{
		Email:   reset.Email,
		Token:   reset.Token,
		Expired: reset.Expired,
	}
}
