package service

import (
	"apriori/entity"
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type NotificationServiceImpl struct {
	NotificationRepository repository.NotificationRepository
	UserRepository         repository.UserRepository
	Notification           entity.NotificationRelation
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

func (service *NotificationServiceImpl) FindAll(ctx context.Context) ([]*model.GetNotificationRelationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	notifications, err := service.NotificationRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	var notificationResponses []*model.GetNotificationRelationResponse
	for _, notification := range notifications {
		notificationResponses = append(notificationResponses, notification.ToNotificationRelationResponse())
	}

	return notificationResponses, nil
}

func (service *NotificationServiceImpl) FindAllByUserId(ctx context.Context, userId int) ([]*model.GetNotificationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	notifications, err := service.NotificationRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	var notificationResponses []*model.GetNotificationResponse
	for _, notification := range notifications {
		notificationResponses = append(notificationResponses, notification.ToNotificationResponse())
	}

	return notificationResponses, nil
}

func (service *NotificationServiceImpl) Create(ctx context.Context, request *model.CreateNotificationRequest) *NotificationServiceImpl {
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
		UserId: request.UserId,
		Title:  strings.Title(request.Title),
		Description: sql.NullString{
			String: request.Description,
			Valid:  true,
		},
		URL: sql.NullString{
			String: request.URL,
			Valid:  true,
		},
		IsRead:    false,
		CreatedAt: timeNow,
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

	service.Notification = entity.NotificationRelation{
		Notification: notificationRequest,
		Name:         user.Name,
		Email:        user.Email,
	}
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
		service.Notification.Email,
		service.Notification.Notification.Title,
		service.Notification.Notification.Description.String,
	)
	if err != nil || service.Error != nil {
		errors := fmt.Errorf("error: %s %s", err.Error(), service.Error.Error())
		return errors
	}

	return nil
}
