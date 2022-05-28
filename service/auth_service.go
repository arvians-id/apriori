package service

import (
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"database/sql"
)

type AuthService interface {
	VerifyCredential(ctx context.Context, request model.GetUserCredentialRequest) (model.GetUserResponse, error)
}

type authService struct {
	UserRepository repository.UserRepository
	AuthRepository repository.AuthRepository
	DB             *sql.DB
}

func NewAuthService(userRepository *repository.UserRepository, authRepository *repository.AuthRepository, db *sql.DB) AuthService {
	return &authService{
		UserRepository: *userRepository,
		AuthRepository: *authRepository,
		DB:             db,
	}
}

func (service *authService) VerifyCredential(ctx context.Context, request model.GetUserCredentialRequest) (model.GetUserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetUserResponse{}, err
	}

	user, err := service.AuthRepository.VerifyCredential(ctx, tx, request.Email, request.Password)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	return utils.ToUserResponse(user), nil
}
