package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"weather-notification/internal/domain/entity"
	handler "weather-notification/internal/domain/error_handler"
	"weather-notification/internal/domain/repository"

	"github.com/google/uuid"
)

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) repository.NotificationRepository {
	return &notificationRepository{
		db: db,
	}
}

func (r *notificationRepository) Create(ctx context.Context, notification *entity.Notification) error {
	query := `
        INSERT INTO notifications (
            id, user_id, location_id, content, status, 
            scheduled_for, sent_at, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	content, err := json.Marshal(notification.Content)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query,
		notification.ID,
		notification.UserID,
		notification.LocationID,
		content,
		notification.Status,
		notification.ScheduledFor,
		notification.SentAt,
		notification.CreatedAt,
		notification.UpdatedAt,
	)

	return err
}

func (r *notificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	query := `
        SELECT id, user_id, location_id, content, status, 
               scheduled_for, sent_at, created_at, updated_at
        FROM notifications
        WHERE id = $1
    `

	notification := &entity.Notification{}
	var content []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.LocationID,
		&content,
		&notification.Status,
		&notification.ScheduledFor,
		&notification.SentAt,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, handler.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &notification.Content)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (r *notificationRepository) FindPendingNotifications(ctx context.Context) ([]*entity.Notification, error) {
	query := `
        SELECT id, user_id, location_id, content, status, 
               scheduled_for, sent_at, created_at, updated_at
        FROM notifications
        WHERE status = $1 AND scheduled_for <= NOW() AT TIME ZONE 'America/Sao_Paulo' 
        ORDER BY scheduled_for
    `

	rows, err := r.db.QueryContext(ctx, query, entity.StatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*entity.Notification

	for rows.Next() {
		notification := &entity.Notification{}
		var content []byte

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.LocationID,
			&content,
			&notification.Status,
			&notification.ScheduledFor,
			&notification.SentAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(content, &notification.Content)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *notificationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.NotificationStatus) error {
	query := `
		UPDATE notifications
		SET status = $1::text, 
			sent_at = CASE 
				WHEN $1::text = 'ENVIADA' THEN NOW() AT TIME ZONE 'America/Sao_Paulo' 
				ELSE sent_at 
			END,
			updated_at = NOW() AT TIME ZONE 'America/Sao_Paulo'
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return handler.ErrNotFound
	}

	return nil
}

func (r *notificationRepository) FindByUserAndLocation(ctx context.Context, userID, locationID uuid.UUID) ([]*entity.Notification, error) {
	query := `
        SELECT id, user_id, location_id, content, status, 
               scheduled_for, sent_at, created_at, updated_at
        FROM notifications
        WHERE user_id = $1 AND location_id = $2
        ORDER BY scheduled_for DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*entity.Notification

	for rows.Next() {
		notification := &entity.Notification{}
		var content []byte

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.LocationID,
			&content,
			&notification.Status,
			&notification.ScheduledFor,
			&notification.SentAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(content, &notification.Content)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *notificationRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	query := `
        SELECT id, user_id, location_id, content, status, 
               scheduled_for, sent_at, created_at, updated_at
        FROM notifications
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*entity.Notification

	for rows.Next() {
		notification := &entity.Notification{}
		var content []byte

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.LocationID,
			&content,
			&notification.Status,
			&notification.ScheduledFor,
			&notification.SentAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(content, &notification.Content)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}
