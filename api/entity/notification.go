package entity

import (
	"time"
)

type Notification struct {
	IdNotification int       `json:"id_notification"`
	UserId         int       `json:"user_id"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	URL            *string   `json:"url"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
	User           *User     `json:"user"`
}
