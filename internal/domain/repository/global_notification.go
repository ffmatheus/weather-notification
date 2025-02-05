package repository

import (
	"context"
	"time"
	"weather-notification/internal/domain/entity"

	"github.com/google/uuid"
)

type GlobalNotificationRepository interface {
	Create(ctx context.Context, notification *entity.GlobalNotification) error
	FindActive(ctx context.Context) ([]*entity.GlobalNotification, error)
	UpdateLastExecution(ctx context.Context, id uuid.UUID, executionTime time.Time) error
}
