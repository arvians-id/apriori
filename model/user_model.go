package model

import "time"

type CreateUserRequest struct {
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateUserRequest struct {
	IdUser    uint64
	Name      string
	Email     string
	Password  string
	UpdatedAt time.Time
}

type GetUserResponse struct {
	IdUser    uint64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
