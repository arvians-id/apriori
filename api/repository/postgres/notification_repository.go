package postgres

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/model"
	"github.com/arvians-id/apriori/repository"
	"log"
)

type NotificationRepositoryImpl struct {
}

func NewNotificationRepository() repository.NotificationRepository {
	return &NotificationRepositoryImpl{}
}

func (repository *NotificationRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Notification, error) {
	query := `SELECT n.*, u.name, u.email FROM notifications n LEFT JOIN users u ON u.id_user = n.user_id ORDER BY n.created_at DESC`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(rows)

	var notifications []*model.Notification
	for rows.Next() {
		notification := model.Notification{
			User: &model.User{},
		}
		err = rows.Scan(
			&notification.IdNotification,
			&notification.UserId,
			&notification.Title,
			&notification.Description,
			&notification.URL,
			&notification.IsRead,
			&notification.CreatedAt,
			&notification.User.Name,
			&notification.User.Email,
		)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

func (repository *NotificationRepositoryImpl) FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*model.Notification, error) {
	query := `SELECT * FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(rows)

	var notifications []*model.Notification
	for rows.Next() {
		var notification model.Notification
		err = rows.Scan(
			&notification.IdNotification,
			&notification.UserId,
			&notification.Title,
			&notification.Description,
			&notification.URL,
			&notification.IsRead,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

func (repository *NotificationRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, notification *model.Notification) (*model.Notification, error) {
	id := 0
	query := `INSERT INTO notifications (user_id, title, description, url, is_read, created_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id_notification`
	row := tx.QueryRowContext(
		ctx,
		query,
		notification.UserId,
		notification.Title,
		notification.Description,
		notification.URL,
		notification.IsRead,
		notification.CreatedAt,
	)
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	notification.IdNotification = id

	return notification, nil
}

func (repository *NotificationRepositoryImpl) Mark(ctx context.Context, tx *sql.Tx, id int) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id_notification = $1`
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (repository *NotificationRepositoryImpl) MarkAll(ctx context.Context, tx *sql.Tx, userId int) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE user_id = $1`
	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}
