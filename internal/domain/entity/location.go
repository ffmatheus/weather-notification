package entity

import (
	"fmt"
	handler "weather-notification/internal/domain/error_handler"

	"github.com/google/uuid"
)

type Location struct {
	ID        uuid.UUID `json:"id"`
	CPTECCode int       `json:"cptecCode"`
	Name      string    `json:"name"`
	State     string    `json:"state"`
}

func NewLocation(cptecCode int, name, state string) (*Location, error) {
	if cptecCode <= 0 {
		return nil, fmt.Errorf("%w: %d", handler.ErrInvalidCPTECCode, cptecCode)
	}

	if name == "" {
		return nil, handler.ErrEmptyLocationName
	}

	if len(state) != 2 {
		return nil, handler.ErrInvalidState
	}

	return &Location{
		ID:        uuid.New(),
		CPTECCode: cptecCode,
		Name:      name,
		State:     state,
	}, nil
}
