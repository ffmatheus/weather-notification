package service

import (
	"context"
	"time"
	"weather-notification/internal/domain/entity"
	"weather-notification/internal/domain/repository"

	"github.com/google/uuid"
)

type GlobalNotificationService struct {
	repo             repository.GlobalNotificationRepository
	userRepo         repository.UserRepository
	queueService     QueueService
	weatherService   *WeatherService
	notificationRepo repository.NotificationRepository
}

func NewGlobalNotificationService(
	repo repository.GlobalNotificationRepository,
	userRepo repository.UserRepository,
	queueService QueueService,
	weatherService *WeatherService,
	notificationRepo repository.NotificationRepository,
) *GlobalNotificationService {
	return &GlobalNotificationService{
		repo:             repo,
		userRepo:         userRepo,
		queueService:     queueService,
		weatherService:   weatherService,
		notificationRepo: notificationRepo,
	}
}

func (s *GlobalNotificationService) Create(ctx context.Context, timeOfDay time.Time, frequency entity.Frequency) error {
	globalNotification := &entity.GlobalNotification{
		ID:        uuid.New(),
		TimeOfDay: timeOfDay,
		Frequency: frequency,
		Active:    true,
		CreatedAt: time.Now(),
	}

	return s.repo.Create(ctx, globalNotification)
}

func (s *GlobalNotificationService) ListActive(ctx context.Context) ([]*entity.GlobalNotification, error) {
	return s.repo.FindActive(ctx)
}

func (s *GlobalNotificationService) ProcessActiveNotifications(ctx context.Context) error {
	now := time.Now()

	notifications, err := s.repo.FindActive(ctx)
	if err != nil {
		return err
	}

	for _, globalNotif := range notifications {
		if !globalNotif.ShouldExecute(now) {
			continue
		}
		users, err := s.userRepo.FindAllActive(ctx)
		if err != nil {
			continue
		}

		for _, user := range users {
			forecast, err := s.weatherService.GetForecast(ctx, user.LocationID)
			if err != nil {
				continue
			}

			notification := &entity.Notification{
				ID:           uuid.New(),
				UserID:       user.ID,
				LocationID:   user.LocationID,
				Content:      *forecast,
				Status:       entity.StatusPending,
				ScheduledFor: now.Add(2 * time.Minute),
				CreatedAt:    now,
			}

			if err := s.notificationRepo.Create(ctx, notification); err != nil {
				return err
			}

			if err := s.queueService.PublishNotification(ctx, notification); err != nil {
				continue
			}
		}

		if err := s.repo.UpdateLastExecution(ctx, globalNotif.ID, now); err != nil {
			continue
		}
	}

	return nil
}
