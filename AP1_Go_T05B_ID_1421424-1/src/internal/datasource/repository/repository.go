package repository

import (
	"context"

	"github.com/google/uuid"

	"tic-tac-toe/internal/domain/model"
)

type GameRepository interface {
	Save(ctx context.Context, game model.Game) error
	Get(ctx context.Context, id uuid.UUID) (model.Game, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListWaiting(ctx context.Context) ([]model.Game, error) // ← новый
	ListFinishedByUser(ctx context.Context, userID uuid.UUID) ([]model.Game, error)
	GetTopPlayers(ctx context.Context, limit int) ([]model.LeaderboardEntry, error)
}

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error) // ← новый
}
