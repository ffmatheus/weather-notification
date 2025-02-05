package entity_test

import (
	"testing"
	"weather-notification/internal/domain/entity"
	handler "weather-notification/internal/domain/error_handler"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	validName := "Matheus"
	validEmail := "matheus@exemplo.com"
	validLocationID := uuid.New()

	tests := []struct {
		name        string
		userName    string
		userEmail   string
		locationID  uuid.UUID
		expectError error
	}{
		{
			name:        "usuário válido",
			userName:    validName,
			userEmail:   validEmail,
			locationID:  validLocationID,
			expectError: nil,
		},
		{
			name:        "nome vazio",
			userName:    "",
			userEmail:   validEmail,
			locationID:  validLocationID,
			expectError: handler.ErrEmptyName,
		},
		{
			name:        "email vazio",
			userName:    validName,
			userEmail:   "",
			locationID:  validLocationID,
			expectError: handler.ErrEmptyEmail,
		},
		{
			name:        "localização inválida",
			userName:    validName,
			userEmail:   validEmail,
			locationID:  uuid.Nil,
			expectError: handler.ErrInvalidLocationID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entity.NewUser(tt.userName, tt.userEmail, tt.locationID)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userName, user.Name)
				assert.Equal(t, tt.userEmail, user.Email)
				assert.Equal(t, tt.locationID, user.LocationID)
				assert.False(t, user.OptOut)
				assert.NotZero(t, user.CreatedAt)
				assert.NotZero(t, user.UpdatedAt)
			}
		})
	}
}
