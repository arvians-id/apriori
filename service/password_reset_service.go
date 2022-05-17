package service

import (
	"apriori/entity"
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
)

type PasswordResetService interface {
	Create(ctx context.Context, request model.CreatePasswordResetRequest, token string, expired int32) (model.GetPasswordResetResponse, error)
}

type passwordResetService struct {
	PasswordResetRepository repository.PasswordResetRepository
	DB                      *sql.DB
}

func NewPasswordResetService(resetRepository repository.PasswordResetRepository, db *sql.DB) PasswordResetService {
	return &passwordResetService{
		PasswordResetRepository: resetRepository,
		DB:                      db,
	}
}

func (service *passwordResetService) Create(ctx context.Context, request model.CreatePasswordResetRequest, token string, expired int32) (model.GetPasswordResetResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetPasswordResetResponse{}, err
	}

	passwordReset := entity.PasswordReset{
		Email:   request.Email,
		Token:   token,
		Expired: expired,
	}

	result, err := service.PasswordResetRepository.Create(ctx, tx, passwordReset)
	if err != nil {
		return model.GetPasswordResetResponse{}, err
	}

	return helper.ToPasswordResetResponse(result), nil
}
