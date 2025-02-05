package postgres

import (
	"context"
	"database/sql"
	"time"
	"weather-notification/internal/domain/entity"
	handler "weather-notification/internal/domain/error_handler"
	"weather-notification/internal/domain/repository"

	"github.com/google/uuid"
)

type globalNotificationRepository struct {
	db *sql.DB
}

func NewGlobalNotificationRepository(db *sql.DB) repository.GlobalNotificationRepository {
	return &globalNotificationRepository{
		db: db,
	}
}

func (r *globalNotificationRepository) Create(ctx context.Context, notification *entity.GlobalNotification) error {
	query := `
        INSERT INTO global_notifications (id, time_of_day, frequency, active, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `

	_, err := r.db.ExecContext(ctx, query,
		notification.ID,
		notification.TimeOfDay,
		notification.Frequency,
		notification.Active,
		notification.CreatedAt,
	)

	return err
}

func (r *globalNotificationRepository) FindActive(ctx context.Context) ([]*entity.GlobalNotification, error) {
	query := `
        SELECT id, time_of_day, frequency, active, last_execution, created_at
        FROM global_notifications
        WHERE active = true
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*entity.GlobalNotification
	for rows.Next() {
		n := &entity.GlobalNotification{}
		err := rows.Scan(
			&n.ID,
			&n.TimeOfDay,
			&n.Frequency,
			&n.Active,
			&n.LastExecution,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *globalNotificationRepository) UpdateLastExecution(ctx context.Context, id uuid.UUID, executionTime time.Time) error {
	query := `
        UPDATE global_notifications
        SET last_execution = $1
        WHERE id = $2
    `

	result, err := r.db.ExecContext(ctx, query, executionTime, id)
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
