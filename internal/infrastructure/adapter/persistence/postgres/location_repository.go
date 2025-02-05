package postgres

import (
	"context"
	"database/sql"
	"weather-notification/internal/domain/entity"
	handler "weather-notification/internal/domain/error_handler"
	"weather-notification/internal/domain/repository"

	"github.com/google/uuid"
)

type locationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) repository.LocationRepository {
	return &locationRepository{
		db: db,
	}
}

func (r *locationRepository) Create(ctx context.Context, location *entity.Location) error {
	query := `
        INSERT INTO locations (id, cptec_id, name, state)
        VALUES ($1, $2, $3, $4)
    `

	_, err := r.db.ExecContext(ctx, query,
		location.ID,
		location.CPTECCode,
		location.Name,
		location.State,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *locationRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Location, error) {
	query := `
        SELECT id, cptec_id, name, state
        FROM locations
        WHERE id = $1
    `

	location := &entity.Location{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&location.ID,
		&location.CPTECCode,
		&location.Name,
		&location.State,
	)

	if err == sql.ErrNoRows {
		return nil, handler.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (r *locationRepository) FindByCPTECCode(ctx context.Context, cptecCode int) (*entity.Location, error) {
	query := `
        SELECT id, cptec_id, name, state
        FROM locations
        WHERE cptec_id = $1
    `

	location := &entity.Location{}
	err := r.db.QueryRowContext(ctx, query, cptecCode).Scan(
		&location.ID,
		&location.CPTECCode,
		&location.Name,
		&location.State,
	)

	if err == sql.ErrNoRows {
		return nil, handler.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (r *locationRepository) FindByNameAndState(ctx context.Context, name, state string) (*entity.Location, error) {
	query := `
        SELECT id, cptec_id, name, state
        FROM locations
        WHERE name = $1 AND state = $2
    `

	location := &entity.Location{}
	err := r.db.QueryRowContext(ctx, query, name, state).Scan(
		&location.ID,
		&location.CPTECCode,
		&location.Name,
		&location.State,
	)

	if err == sql.ErrNoRows {
		return nil, handler.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return location, nil
}
