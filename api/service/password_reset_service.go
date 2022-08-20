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
	CreateOrUpdateByEmail(ctx context.Context, email string) (model.GetPasswordResetResponse, error)
	Verify(ctx context.Context, request model.UpdateResetPasswordUserRequest) error
}

type passwordResetService struct {
	PasswordResetRepository repository.PasswordResetRepository
	UserRepository          repository.UserRepository
	DB                      *sql.DB
}

func NewPasswordResetService(resetRepository *repository.PasswordResetRepository, userRepository *repository.UserRepository, db *sql.DB) PasswordResetService {
	return &passwordResetService{
		PasswordResetRepository: *resetRepository,
		UserRepository:          *userRepository,
		DB:                      db,
	}
}

func (service *passwordResetService) CreateOrUpdateByEmail(ctx context.Context, email string) (model.GetPasswordResetResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetPasswordResetResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	timestamp := time.Now().Add(1 * time.Hour).Unix()
	timestampString := strconv.Itoa(int(timestamp))
	token := md5.Sum([]byte(email + timestampString))
	tokenString := fmt.Sprintf("%x", token)
	passwordResetRequest := entity.PasswordReset{
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
	_, err = service.PasswordResetRepository.FindByEmail(ctx, tx, user.Email)
	if err != nil {
		// Create new data if not exists
		passwordResetResponse, err := service.PasswordResetRepository.Create(ctx, tx, passwordResetRequest)
		if err != nil {
			return model.GetPasswordResetResponse{}, err
		}

		return utils.ToPasswordResetResponse(passwordResetResponse), nil
	}

	// Update data if exists
	passwordResetResponse, err := service.PasswordResetRepository.Update(ctx, tx, passwordResetRequest)
	if err != nil {
		return model.GetPasswordResetResponse{}, err
	}

	return utils.ToPasswordResetResponse(passwordResetResponse), nil
}

func (service *passwordResetService) Verify(ctx context.Context, request model.UpdateResetPasswordUserRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	// Check if email and token is exists in table password_resets
	passwordResetRequest := entity.PasswordReset{
		Email: request.Email,
		Token: request.Token,
	}

	reset, err := service.PasswordResetRepository.FindByEmailAndToken(ctx, tx, passwordResetRequest)
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
	timeNow, err := time.Parse(utils.TimeFormat, now.Format(utils.TimeFormat))
	if err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(password)
	user.UpdatedAt = timeNow

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
