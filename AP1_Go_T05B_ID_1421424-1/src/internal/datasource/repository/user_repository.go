package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"tic-tac-toe/internal/domain/model"
	apperrors "tic-tac-toe/internal/errors"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, u model.User) error {
	const q = `INSERT INTO users (id, login, password_hash) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, q, u.ID, u.Login, u.PasswordHash)
	if err != nil {
		// нарушение UNIQUE по login
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	const q = `SELECT id, login, password_hash FROM users WHERE login = $1`

	var u model.User
	err := r.pool.QueryRow(ctx, q, login).Scan(&u.ID, &u.Login, &u.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, apperrors.ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("get user by login: %w", err)
	}
	return u, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	const q = `SELECT id, login, password_hash FROM users WHERE id = $1`

	var u model.User
	err := r.pool.QueryRow(ctx, q, id).Scan(&u.ID, &u.Login, &u.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, apperrors.ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}
