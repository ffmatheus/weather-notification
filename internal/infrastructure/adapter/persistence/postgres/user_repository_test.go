package postgres_test

import (
	"context"
	"testing"
	"time"
	"weather-notification/internal/domain/entity"
	"weather-notification/internal/infrastructure/adapter/persistence/postgres"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro criando mock do db: %v", err)
	}
	defer db.Close()

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{
		ID:         uuid.New(),
		LocationID: uuid.New(),
		Name:       "Matheus",
		Email:      "matheus@exemplo.com",
		OptOut:     false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mock.ExpectExec(`INSERT INTO users`).
		WithArgs(user.ID, user.LocationID, user.Name, user.Email, user.OptOut, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
