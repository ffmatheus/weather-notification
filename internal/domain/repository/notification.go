package repository

import (
	"context"
	"weather-notification/internal/domain/entity" // ajuste para seu path

	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	FindPendingNotifications(ctx context.Context) ([]*entity.Notification, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.NotificationStatus) error
	FindByUserAndLocation(ctx context.Context, userID, locationID uuid.UUID) ([]*entity.Notification, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error)
}
