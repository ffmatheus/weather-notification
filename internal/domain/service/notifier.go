package service

import (
	"context"
	"weather-notification/internal/domain/entity"
)

type Notifier interface {
	Send(ctx context.Context, notification *entity.Notification) error
}
