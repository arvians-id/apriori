package request

type CreateNotificationRequest struct {
	UserId      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	IsRead      bool   `json:"is_read"`
	CreatedAt   string `json:"created_at"`
}
