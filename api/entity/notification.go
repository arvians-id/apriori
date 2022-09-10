package entity

import (
	"apriori/model"
	"database/sql"
	"time"
)

type Notification struct {
	IdNotification int
	UserId         int
	Title          string
	Description    sql.NullString
	URL            sql.NullString
	IsRead         bool
	CreatedAt      time.Time
}

func (notification *Notification) ToNotificationResponse() *model.GetNotificationResponse {
	return &model.GetNotificationResponse{
		IdNotification: notification.IdNotification,
		UserId:         notification.UserId,
		Title:          notification.Title,
		Description:    notification.Description.String,
		URL:            notification.URL.String,
		IsRead:         notification.IsRead,
		CreatedAt:      notification.CreatedAt.String(),
	}
}

type NotificationRelation struct {
	Notification Notification
	Name         string
	Email        string
}

func (notification *NotificationRelation) ToNotificationRelationResponse() *model.GetNotificationRelationResponse {
	return &model.GetNotificationRelationResponse{
		IdNotification: notification.Notification.IdNotification,
		UserId:         notification.Notification.UserId,
		Title:          notification.Notification.Title,
		Description:    notification.Notification.Description.String,
		URL:            notification.Notification.URL.String,
		IsRead:         notification.Notification.IsRead,
		CreatedAt:      notification.Notification.CreatedAt.String(),
		Name:           notification.Name,
		Email:          notification.Email,
	}
}
