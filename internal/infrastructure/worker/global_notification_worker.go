package worker

import (
	"context"
	"log"
	"time"
	"weather-notification/internal/domain/service"
)

type GlobalNotificationWorker struct {
	ctx     context.Context
	service *service.GlobalNotificationService
}

func NewGlobalNotificationWorker(
	ctx context.Context,
	service *service.GlobalNotificationService,
) *GlobalNotificationWorker {
	return &GlobalNotificationWorker{
		ctx:     ctx,
		service: service,
	}
}

func (w *GlobalNotificationWorker) Start() error {
	log.Printf("Iniciando worker de notificações globais...")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			log.Printf("Finalizando worker de notificações globais...")
			return nil
		case <-ticker.C:
			if err := w.service.ProcessActiveNotifications(w.ctx); err != nil {
				log.Printf("Erro ao processar notificações globais: %v", err)
			}
		}
	}
}
