package service

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"github.com/arvians-id/apriori/util"
	"log"
	"strings"
	"time"
)

type NotificationServiceImpl struct {
	NotificationRepository repository.NotificationRepository
	UserRepository         repository.UserRepository
	Notification           *model.Notification
	Error                  error
	DB                     *sql.DB
}

func NewNotificationService(
	notificationRepository *repository.NotificationRepository,
	userRepository *repository.UserRepository,
	db *sql.DB,
) NotificationService {
	return &NotificationServiceImpl{
		NotificationRepository: *notificationRepository,
		UserRepository:         *userRepository,
		DB:                     db,
	}
}

func (service *NotificationServiceImpl) FindAll(ctx context.Context) ([]*model.Notification, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	notifications, err := service.NotificationRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[NotificationService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return notifications, nil
}

func (service *NotificationServiceImpl) FindAllByUserId(ctx context.Context, userId int) ([]*model.Notification, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][FindAllByUserId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	notifications, err := service.NotificationRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		log.Println("[NotificationService][FindAllByUserId][FindAllByUserId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return notifications, nil
}

func (service *NotificationServiceImpl) Create(ctx context.Context, request *request.CreateNotificationRequest) (*model.Notification, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[NotificationService][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	notificationRequest := model.Notification{
		UserId:      request.UserId,
		Title:       strings.Title(request.Title),
		Description: &request.Description,
		URL:         &request.URL,
		IsRead:      false,
		CreatedAt:   timeNow,
	}

	notificationResponse, err := service.NotificationRepository.Create(ctx, tx, &notificationRequest)
	if err != nil {
		log.Println("[NotificationService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return notificationResponse, nil
}

func (service *NotificationServiceImpl) MarkAll(ctx context.Context, userId int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][MarkAll] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	err = service.NotificationRepository.MarkAll(ctx, tx, userId)
	if err != nil {
		log.Println("[NotificationService][Create][MarkAll] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}

func (service *NotificationServiceImpl) Mark(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][Mark] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	err = service.NotificationRepository.Mark(ctx, tx, id)
	if err != nil {
		log.Println("[NotificationService][Create][Mark] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}
