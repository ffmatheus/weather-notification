package entity

import (
	"time"
	handler "weather-notification/internal/domain/error_handler"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID
	LocationID uuid.UUID
	Name       string
	Email      string
	OptOut     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewUser(name, email string, locationID uuid.UUID) (*User, error) {
	if name == "" {
		return nil, handler.ErrEmptyName
	}
	if email == "" {
		return nil, handler.ErrEmptyEmail
	}
	if locationID == uuid.Nil {
		return nil, handler.ErrInvalidLocationID
	}

	now := time.Now()
	return &User{
		ID:         uuid.New(),
		Name:       name,
		Email:      email,
		LocationID: locationID,
		OptOut:     false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
