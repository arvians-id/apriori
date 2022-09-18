package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arvians-id/apriori/entity"
	"github.com/arvians-id/apriori/helper"
	"github.com/arvians-id/apriori/http/request"
	"github.com/arvians-id/apriori/repository"
	"strings"
	"time"
)

type NotificationServiceImpl struct {
	NotificationRepository repository.NotificationRepository
	UserRepository         repository.UserRepository
	Notification           *entity.Notification
	EmailService           EmailService
	Error                  error
	DB                     *sql.DB
}

func NewNotificationService(notificationRepository *repository.NotificationRepository, userRepository *repository.UserRepository, emailService *EmailService, db *sql.DB) NotificationService {
	return &NotificationServiceImpl{
		NotificationRepository: *notificationRepository,
		UserRepository:         *userRepository,
		EmailService:           *emailService,
		DB:                     db,
	}
}

func (service *NotificationServiceImpl) FindAll(ctx context.Context) ([]*entity.Notification, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	notifications, err := service.NotificationRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (service *NotificationServiceImpl) FindAllByUserId(ctx context.Context, userId int) ([]*entity.Notification, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	notifications, err := service.NotificationRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (service *NotificationServiceImpl) Create(ctx context.Context, request *request.CreateNotificationRequest) *NotificationServiceImpl {
	tx, err := service.DB.Begin()
	if err != nil {
		service.Error = err
		return nil
	}
	defer helper.CommitOrRollback(tx)

	timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
	if err != nil {
		service.Error = err
		return nil
	}

	notificationRequest := entity.Notification{
		UserId:      request.UserId,
		Title:       strings.Title(request.Title),
		Description: &request.Description,
		URL:         &request.URL,
		IsRead:      false,
		CreatedAt:   timeNow,
	}

	notificationResponse, err := service.NotificationRepository.Create(ctx, tx, &notificationRequest)
	if err != nil {
		service.Error = err
		return nil
	}

	user, err := service.UserRepository.FindById(ctx, tx, notificationResponse.UserId)
	if err != nil {
		service.Error = err
		return nil
	}

	notificationRequest.User = &entity.User{
		Name:  user.Name,
		Email: user.Email,
	}

	service.Notification = &notificationRequest
	return service
}

func (service *NotificationServiceImpl) MarkAll(ctx context.Context, userId int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.NotificationRepository.MarkAll(ctx, tx, userId)
	if err != nil {
		return err
	}

	return nil
}

func (service *NotificationServiceImpl) Mark(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.NotificationRepository.Mark(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *NotificationServiceImpl) WithSendMail() error {
	err := service.EmailService.SendEmailWithText(
		service.Notification.User.Email,
		service.Notification.Title,
		service.Notification.Description,
	)
	if err != nil || service.Error != nil {
		errors := fmt.Errorf("error: %s %s", err.Error(), service.Error.Error())
		return errors
	}

	return nil
}
