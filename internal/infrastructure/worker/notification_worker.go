package worker

import (
	"context"
	"log"
	"os"
	"time"
	"weather-notification/internal/domain/entity"
	"weather-notification/internal/domain/service"
	"weather-notification/internal/infrastructure/adapter/notifier"
)

type NotificationWorker struct {
	ctx             context.Context
	notificationSvc *service.NotificationService
	weatherSvc      *service.WeatherService
	queueService    service.QueueService
}

func NewNotificationWorker(
	ctx context.Context,
	notificationSvc *service.NotificationService,
	weatherSvc *service.WeatherService,
	queueService service.QueueService,
) *NotificationWorker {
	return &NotificationWorker{
		ctx:             ctx,
		notificationSvc: notificationSvc,
		weatherSvc:      weatherSvc,
		queueService:    queueService,
	}
}

func (w *NotificationWorker) Start() error {
	log.Printf("Iniciando worker de notificações...")
	return w.queueService.ConsumeNotifications(w.ctx, w.processNotification)
}

func (w *NotificationWorker) processNotification(notification *entity.Notification) error {
	log.Printf("Recebendo notificação para processamento: %s", notification.ID)

	if time.Now().Before(notification.ScheduledFor) {
		log.Printf("Notificação %s agendada para futuro: %v", notification.ID, notification.ScheduledFor)
		return nil
	}

	forecast, err := w.weatherSvc.GetForecast(w.ctx, notification.LocationID)
	if err != nil {
		log.Printf("Erro ao atualizar previsão: %v", err)
		return err
	}

	notification.Content = *forecast

	err = w.sendNotification(notification)
	if err != nil {
		log.Printf("Erro ao enviar notificação: %v", err)
		return err
	}

	notification.MarkAsSent()
	err = w.notificationSvc.UpdateStatus(w.ctx, notification)
	if err != nil {
		log.Printf("Erro ao atualizar status: %v", err)
		return err
	}

	return nil
}

func (w *NotificationWorker) sendNotification(notification *entity.Notification) error {
	notifier := notifier.NewWebNotifier(os.Getenv("WEBHOOK_URL"))

	err := notifier.Send(w.ctx, notification)
	if err != nil {
		log.Printf("Erro ao enviar notificação web: %v", err)
		return err
	}

	return nil
}
