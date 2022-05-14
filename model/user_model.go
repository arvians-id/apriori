package model

type CreateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateUserRequest struct {
	IdUser    uint64 `json:"id_user"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	UpdatedAt string `json:"updated_at"`
}

type GetUserResponse struct {
	IdUser    uint64 `json:"id_user"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
