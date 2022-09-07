package entity

import (
	"apriori/model"
	"time"
)

type User struct {
	IdUser    int
	Role      int
	Name      string
	Email     string
	Address   string
	Phone     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *User) ToUserResponse() *model.GetUserResponse {
	return &model.GetUserResponse{
		IdUser:    user.IdUser,
		Role:      user.Role,
		Name:      user.Name,
		Email:     user.Email,
		Address:   user.Address,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}
