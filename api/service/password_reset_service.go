package service

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"github.com/arvians-id/apriori/helper"
	"github.com/arvians-id/apriori/http/request"
	"github.com/arvians-id/apriori/model"
	"github.com/arvians-id/apriori/repository"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type PasswordResetServiceImpl struct {
	PasswordResetRepository repository.PasswordResetRepository
	UserRepository          repository.UserRepository
	DB                      *sql.DB
}

func NewPasswordResetService(
	resetRepository *repository.PasswordResetRepository,
	userRepository *repository.UserRepository,
	db *sql.DB,
) PasswordResetService {
	return &PasswordResetServiceImpl{
		PasswordResetRepository: *resetRepository,
		UserRepository:          *userRepository,
		DB:                      db,
	}
}

func (service *PasswordResetServiceImpl) CreateOrUpdateByEmail(ctx context.Context, email string) (*model.PasswordReset, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	timestamp := time.Now().Add(1 * time.Hour).Unix()
	timestampString := strconv.Itoa(int(timestamp))
	token := md5.Sum([]byte(email + timestampString))
	tokenString := fmt.Sprintf("%x", token)
	passwordResetRequest := model.PasswordReset{
		Email:   email,
		Token:   tokenString,
		Expired: timestamp,
	}

	// Check if email is exists in table users
	user, err := service.UserRepository.FindByEmail(ctx, tx, email)
	if err != nil {
		return nil, err
	}

	// Check If email is exists in table password_resets
	_, err = service.PasswordResetRepository.FindByEmail(ctx, tx, user.Email)
	if err != nil {
		// Create new data if not exists
		passwordReset, err := service.PasswordResetRepository.Create(ctx, tx, &passwordResetRequest)
		if err != nil {
			return nil, err
		}

		return passwordReset, nil
	}

	// Update data if exists
	passwordReset, err := service.PasswordResetRepository.Update(ctx, tx, &passwordResetRequest)
	if err != nil {
		return nil, err
	}

	return passwordReset, nil
}

func (service *PasswordResetServiceImpl) Verify(ctx context.Context, request *request.UpdateResetPasswordUserRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Check if email and token is exists in table password_resets
	passwordResetRequest := model.PasswordReset{
		Email: request.Email,
		Token: request.Token,
	}

	reset, err := service.PasswordResetRepository.FindByEmailAndToken(ctx, tx, &passwordResetRequest)
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
	timeNow, err := time.Parse(helper.TimeFormat, now.Format(helper.TimeFormat))
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
