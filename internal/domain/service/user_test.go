package service_test

import (
	"context"
	"testing"
	"weather-notification/internal/domain/entity"
	"weather-notification/internal/domain/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockUserRepository) FindAllActive(ctx context.Context) ([]entity.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateOptOut(ctx context.Context, id uuid.UUID, optOut bool) error {
	args := m.Called(ctx, id, optOut)
	return args.Error(0)
}

func TestUserService_Create(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	ctx := context.Background()

	validName := "Matheus"
	validEmail := "matheus@exemplo.com"
	validLocationID := uuid.New()

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		locationID   uuid.UUID
		mockBehavior func(mockRepo *MockUserRepository)
		expectError  bool
	}{
		{
			name:       "sucesso ao criar usuário",
			userName:   validName,
			userEmail:  validEmail,
			locationID: validLocationID,
			mockBehavior: func(mockRepo *MockUserRepository) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(user *entity.User) bool {
					return user.Name == validName &&
						user.Email == validEmail &&
						user.LocationID == validLocationID
				})).Return(nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRepo)

			err := userService.Create(ctx, tt.userName, tt.userEmail, tt.locationID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
