package repository

import (
	"context"
	"weather-notification/internal/domain/entity"

	"github.com/google/uuid"
)

type LocationRepository interface {
	Create(ctx context.Context, location *entity.Location) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Location, error)
	FindByCPTECCode(ctx context.Context, cptecCode int) (*entity.Location, error)
	FindByNameAndState(ctx context.Context, name, state string) (*entity.Location, error)
}
