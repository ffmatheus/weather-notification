package entity

import (
	"time"
	handler "weather-notification/internal/domain/error_handler"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	LocationID uuid.UUID `json:"location_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	OptOut     bool      `json:"opt_out"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
