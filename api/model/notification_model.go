package model

type GetNotificationResponse struct {
	IdNotification int    `json:"id_notification"`
	UserId         int    `json:"user_id"`
	Title          string `json:"title"`
	Description    string `json:"description,omitempty"`
	URL            string `json:"url,omitempty"`
	IsRead         bool   `json:"is_read"`
	CreatedAt      string `json:"created_at"`
}

type CreateNotificationRequest struct {
	UserId      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	IsRead      bool   `json:"is_read"`
	CreatedAt   string `json:"created_at"`
}

type GetNotificationRelationResponse struct {
	IdNotification int    `json:"id_notification"`
	UserId         int    `json:"user_id"`
	Title          string `json:"title"`
	Description    string `json:"description,omitempty"`
	URL            string `json:"url,omitempty"`
	IsRead         bool   `json:"is_read"`
	CreatedAt      string `json:"created_at"`
	Name           string `json:"name"`
	Email          string `json:"email"`
}
