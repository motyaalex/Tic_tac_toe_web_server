package service

import (
	"context"

	"tic-tac-toe/internal/domain/model"

	"github.com/google/uuid"
)

type GameService interface {
	CreateGame(ctx context.Context, userID uuid.UUID, vsComputer bool) (model.Game, error)
	JoinGame(ctx context.Context, gameID, userID uuid.UUID) (model.Game, error)
	MakeMove(ctx context.Context, gameID, userID uuid.UUID, row, col int) (model.Game, error)
	GetGame(ctx context.Context, gameID uuid.UUID) (model.Game, error)
	ListWaiting(ctx context.Context) ([]model.Game, error)
	DeleteGame(ctx context.Context, id uuid.UUID) error
	ListFinishedByUser(ctx context.Context, userID uuid.UUID) ([]model.Game, error)
	GetTopPlayers(ctx context.Context, limit int) ([]model.LeaderboardEntry, error)
}
