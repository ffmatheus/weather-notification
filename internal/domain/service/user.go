package service

import (
	"context"
	"weather-notification/internal/domain/entity"
	"weather-notification/internal/domain/repository"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context, name, email string, locationID uuid.UUID) error {
	user, _ := entity.NewUser(name, email, locationID)
	return s.userRepo.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, userID uuid.UUID, name string, locationID uuid.UUID) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if name != "" {
		user.Name = name
	}

	if locationID != uuid.Nil {
		user.LocationID = locationID
	}

	return s.userRepo.Update(ctx, user)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) ToggleOptOut(ctx context.Context, userID uuid.UUID, optOut bool) error {
	return s.userRepo.UpdateOptOut(ctx, userID, optOut)
}
