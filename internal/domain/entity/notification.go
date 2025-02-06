package entity

import (
	"time"
	handler "weather-notification/internal/domain/error_handler"

	"github.com/google/uuid"
)

type NotificationStatus string
type Frequency string

const (
	StatusPending   NotificationStatus = "PENDENTE"
	StatusSent      NotificationStatus = "ENVIADA"
	StatusFailed    NotificationStatus = "FALHA"
	FrequencyDaily  Frequency          = "DIARIA"
	FrequencyWeekly Frequency          = "SEMANAL"
)

type GlobalNotification struct {
	ID            uuid.UUID  `json:"id"`
	TimeOfDay     time.Time  `json:"time_of_day"`
	Frequency     Frequency  `json:"frequency"`
	Active        bool       `json:"active"`
	LastExecution *time.Time `json:"last_execution"`
	CreatedAt     time.Time  `json:"created_at"`
}

type Notification struct {
	ID           uuid.UUID                 `json:"id"`
	UserID       uuid.UUID                 `json:"user_id"`
	LocationID   uuid.UUID                 `json:"location_id"`
	Content      WeatherForecastCollection `json:"content"`
	Status       NotificationStatus        `json:"status"`
	ScheduledFor time.Time                 `json:"scheduled_for"`
	SentAt       *time.Time                `json:"sent_at"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

func NewNotification(userID, locationID uuid.UUID, content WeatherForecastCollection, scheduledFor time.Time) (*Notification, error) {
	if userID == uuid.Nil {
		return nil, handler.ErrInvalidUserID
	}
	if locationID == uuid.Nil {
		return nil, handler.ErrInvalidLocationID
	}

	now := time.Now().Truncate(time.Second)
	scheduledFor = scheduledFor.Truncate(time.Second)

	if scheduledFor.Before(now) {
		return nil, handler.ErrInvalidScheduleTime
	}

	return &Notification{
		ID:           uuid.New(),
		UserID:       userID,
		LocationID:   locationID,
		Content:      content,
		Status:       StatusPending,
		ScheduledFor: scheduledFor,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (g *GlobalNotification) ShouldExecute(now time.Time) bool {
	if !g.Active {
		return false
	}

	currentTime := now.Format("15:04")
	scheduledTime := g.TimeOfDay.Format("15:04")

	if currentTime != scheduledTime {
		return false
	}

	if g.LastExecution == nil {
		return true
	}

	switch g.Frequency {
	case FrequencyDaily:
		return g.LastExecution.Day() != now.Day()
	case FrequencyWeekly:
		_, lastWeek := g.LastExecution.ISOWeek()
		_, currentWeek := now.ISOWeek()
		return lastWeek != currentWeek
	default:
		return false
	}
}

func (n *Notification) IsReadyToSend() bool {
	return n.Status == StatusPending && time.Now().After(n.ScheduledFor)
}

func (n *Notification) MarkAsSent() {
	now := time.Now()
	n.Status = StatusSent
	n.SentAt = &now
	n.UpdatedAt = now
}

func (n *Notification) MarkAsFailed() {
	n.Status = StatusFailed
	n.UpdatedAt = time.Now()
}

func (n *Notification) FormatNotificationContent() string {
	forecasts := n.Content.GetNext4Days()
	result := "Previsão do tempo para os próximos dias:\n\n"

	for _, forecast := range forecasts {
		result += forecast.AsNotificationText() + "\n"
	}

	return result
}

func (n *Notification) ValidateForSending() error {
	if n.Status != StatusPending {
		return handler.ErrInvalidNotificationStatus
	}
	if len(n.Content.Forecasts) == 0 {
		return handler.ErrEmptyForecast
	}
	return nil
}
