package service

import (
	"context"
	"weather-notification/internal/domain/entity"
)

type QueueService interface {
	PublishNotification(ctx context.Context, notification *entity.Notification) error
	ConsumeNotifications(ctx context.Context, handler func(*entity.Notification) error) error
	Close() error
}
