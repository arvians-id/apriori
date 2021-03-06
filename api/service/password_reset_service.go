package service

import (
	"apriori/entity"
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type PasswordResetService interface {
	CreateOrUpdate(ctx context.Context, email string) (model.GetPasswordResetResponse, error)
	Verify(ctx context.Context, request model.UpdateResetPasswordUserRequest) error
}

type passwordResetService struct {
	PasswordResetRepository repository.PasswordResetRepository
	UserRepository          repository.UserRepository
	DB                      *sql.DB
	date                    string
}

func NewPasswordResetService(resetRepository *repository.PasswordResetRepository, userRepository *repository.UserRepository, db *sql.DB) PasswordResetService {
	return &passwordResetService{
		PasswordResetRepository: *resetRepository,
		UserRepository:          *userRepository,
		DB:                      db,
		date:                    "2006-01-02 15:04:05",
	}
}

func (service *passwordResetService) CreateOrUpdate(ctx context.Context, email string) (model.GetPasswordResetResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetPasswordResetResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	timestamp := time.Now().Add(1 * time.Hour).Unix()
	timestampString := strconv.Itoa(int(timestamp))

	token := md5.Sum([]byte(email + timestampString))
	tokenString := fmt.Sprintf("%x", token)
	passwordReset := entity.PasswordReset{
		Email:   email,
		Token:   tokenString,
		Expired: int32(timestamp),
	}

	// Check if email is exists in table users
	user, err := service.UserRepository.FindByEmail(ctx, tx, email)
	if err != nil {
		return model.GetPasswordResetResponse{}, err
	}

	// Check If email is exists in table password_resets
	rows, err := service.PasswordResetRepository.FindByEmail(ctx, tx, user.Email)
	if err != nil {
		return model.GetPasswordResetResponse{}, err
	}

	// Create new data if not exists
	if rows.Email == "" {
		result, err := service.PasswordResetRepository.Create(ctx, tx, passwordReset)
		if err != nil {
			return model.GetPasswordResetResponse{}, err
		}

		return utils.ToPasswordResetResponse(result), nil
	}

	// Update data if exists
	result, err := service.PasswordResetRepository.Update(ctx, tx, passwordReset)
	if err != nil {
		return model.GetPasswordResetResponse{}, err
	}

	return utils.ToPasswordResetResponse(result), nil
}

func (service *passwordResetService) Verify(ctx context.Context, request model.UpdateResetPasswordUserRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	// Check if email and token is exists in table password_resets
	passwordReset := entity.PasswordReset{
		Email: request.Email,
		Token: request.Token,
	}

	reset, err := service.PasswordResetRepository.FindByEmailAndToken(ctx, tx, passwordReset)
	if err != nil {
		return err
	}

	// Check token expired
	now := time.Now()

	// if expired
	if now.Unix() > int64(reset.Expired) {
		err := service.PasswordResetRepository.Delete(ctx, tx, reset.Email)
		if err != nil {
			return err
		}

		return errors.New("reset password verification is expired")
	}

	// if not
	// Check if email is exists in table users
	user, err := service.UserRepository.FindByEmail(ctx, tx, reset.Email)
	if err != nil {
		return err
	}

	// Update the password
	updatedAt, err := time.Parse(service.date, now.Format(service.date))
	if err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(password)
	user.UpdatedAt = updatedAt

	err = service.UserRepository.UpdatePassword(ctx, tx, user)
	if err != nil {
		return err
	}

	// Delete data from table password_reset
	err = service.PasswordResetRepository.Delete(ctx, tx, user.Email)
	if err != nil {
		return err
	}

	return nil
}
