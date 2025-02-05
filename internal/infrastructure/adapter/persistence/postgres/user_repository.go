package postgres

import (
	"context"
	"database/sql"
	"weather-notification/internal/domain/entity"
	handler "weather-notification/internal/domain/error_handler"
	"weather-notification/internal/domain/repository"

	"github.com/google/uuid"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
        INSERT INTO users (id, location_id, name, email, opt_out, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.LocationID,
		user.Name,
		user.Email,
		user.OptOut,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, location_id = $3, opt_out = $4, updated_at = NOW()
		WHERE id = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.LocationID,
		user.OptOut,
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
        SELECT id, name, email, opt_out, created_at, updated_at
        FROM users
        WHERE id = $1
    `
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.OptOut,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, handler.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
        SELECT id, name, email, opt_out, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.OptOut,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, handler.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateOptOut(ctx context.Context, id uuid.UUID, optOut bool) error {
	query := `
        UPDATE users
        SET opt_out = $1, updated_at = NOW()
        WHERE id = $2
    `

	result, err := r.db.ExecContext(ctx, query, optOut, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return handler.ErrNotFound
	}

	return nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	query := `
		SELECT id, name, email, opt_out, created_at, updated_at
		FROM users
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		user := entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.OptOut,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) FindAllActive(ctx context.Context) ([]entity.User, error) {
	query := `
		SELECT id, location_id, name, email, opt_out, created_at, updated_at
		FROM users
		WHERE opt_out = false
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		user := entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.LocationID,
			&user.Name,
			&user.Email,
			&user.OptOut,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
