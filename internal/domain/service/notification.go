package service

import (
	"context"
	"time"
	"weather-notification/internal/domain/entity"
	handler "weather-notification/internal/domain/error_handler"
	"weather-notification/internal/domain/repository"

	"github.com/google/uuid"
)

type NotificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
	weatherService   *WeatherService
	queueService     QueueService
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	weatherService *WeatherService,
	queueService QueueService,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		weatherService:   weatherService,
		queueService:     queueService,
	}
}

func (s *NotificationService) Schedule(ctx context.Context, userID, locationID uuid.UUID, scheduledFor time.Time) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.OptOut {
		return handler.ErrUserOptOut
	}

	forecast, err := s.weatherService.GetForecast(ctx, locationID)
	if err != nil {
		return err
	}

	notification, err := entity.NewNotification(
		userID,
		locationID,
		*forecast,
		scheduledFor,
	)
	if err != nil {
		return err
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	return s.queueService.PublishNotification(ctx, notification)
}

func (s *NotificationService) ProcessPendingNotifications(ctx context.Context) error {
	notifications, err := s.notificationRepo.FindPendingNotifications(ctx)
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		user, err := s.userRepo.FindByID(ctx, notification.UserID)
		if err != nil {
			continue
		}
		if user.OptOut {
			notification.MarkAsFailed()
			s.notificationRepo.UpdateStatus(ctx, notification.ID, entity.StatusFailed)
			continue
		}

		forecast, err := s.weatherService.GetForecast(ctx, notification.LocationID)
		if err != nil {
			notification.MarkAsFailed()
			s.notificationRepo.UpdateStatus(ctx, notification.ID, entity.StatusFailed)
			continue
		}

		notification.Content = *forecast

		notification.MarkAsSent()
		s.notificationRepo.UpdateStatus(ctx, notification.ID, entity.StatusSent)
	}

	return nil
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.notificationRepo.FindByUser(ctx, userID)
}

func (s *NotificationService) UpdateStatus(ctx context.Context, notification *entity.Notification) error {
	return s.notificationRepo.UpdateStatus(ctx, notification.ID, notification.Status)
}
