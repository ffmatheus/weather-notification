package repository

import (
	"context"
	"weather-notification/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	FindAll(ctx context.Context) ([]entity.User, error)
	FindAllActive(ctx context.Context) ([]entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateOptOut(ctx context.Context, id uuid.UUID, optOut bool) error
}
