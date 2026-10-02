package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"tic-tac-toe/internal/domain/model"
	apperrors "tic-tac-toe/internal/errors"
)

type gameRepository struct {
	pool *pgxpool.Pool
}

func NewGameRepository(pool *pgxpool.Pool) GameRepository {
	return &gameRepository{pool: pool}
}

func (r *gameRepository) Save(ctx context.Context, game model.Game) error {
	boardJSON, err := json.Marshal(game.Board)
	if err != nil {
		return fmt.Errorf("marshal board: %w", err)
	}

	const q = `
		INSERT INTO games (id, board, player_x, player_o, state, turn_player, winner, vs_computer, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			board       = EXCLUDED.board,
			player_x    = EXCLUDED.player_x,
			player_o    = EXCLUDED.player_o,
			state       = EXCLUDED.state,
			turn_player = EXCLUDED.turn_player,
			winner      = EXCLUDED.winner,
			vs_computer = EXCLUDED.vs_computer
	`

	_, err = r.pool.Exec(ctx, q,
		game.ID,
		boardJSON,
		game.PlayerX,
		game.PlayerO,
		game.State,
		game.TurnPlayer,
		game.Winner,
		game.VsComputer,
		game.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("exec save game: %w", err)
	}
	return nil
}

func (r *gameRepository) Get(ctx context.Context, id uuid.UUID) (model.Game, error) {
	const q = `
		SELECT id, board, player_x, player_o, state, turn_player, winner, vs_computer, created_at
		FROM games
		WHERE id = $1
	`

	g, err := scanGame(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Game{}, apperrors.ErrNotFound
	}
	if err != nil {
		return model.Game{}, fmt.Errorf("get game: %w", err)
	}
	return g, nil
}

func (r *gameRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM games WHERE id = $1`
	if _, err := r.pool.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("delete game: %w", err)
	}
	return nil
}

func (r *gameRepository) ListWaiting(ctx context.Context) ([]model.Game, error) {
	const q = `
		SELECT id, board, player_x, player_o, state, turn_player, winner, vs_computer, created_at
		FROM games
		WHERE state = $1
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, q, model.StateWaitingForPlayers)
	if err != nil {
		return nil, fmt.Errorf("query waiting games: %w", err)
	}
	defer rows.Close()

	var games []model.Game
	for rows.Next() {
		g, err := scanGame(rows)
		if err != nil {
			return nil, fmt.Errorf("scan waiting game: %w", err)
		}
		games = append(games, g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate waiting games: %w", err)
	}
	return games, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

// scanGame читает одну строку games во всех методах одинаково.

func scanGame(row rowScanner) (model.Game, error) {
	var (
		g         model.Game
		boardJSON []byte
	)
	err := row.Scan(
		&g.ID,
		&boardJSON,
		&g.PlayerX,
		&g.PlayerO,
		&g.State,
		&g.TurnPlayer,
		&g.Winner,
		&g.VsComputer,
		&g.CreatedAt,
	)
	if err != nil {
		return model.Game{}, err
	}
	if err := json.Unmarshal(boardJSON, &g.Board); err != nil {
		return model.Game{}, fmt.Errorf("unmarshal board: %w", err)
	}
	return g, nil
}

func (r *gameRepository) ListFinishedByUser(ctx context.Context, userID uuid.UUID) ([]model.Game, error) {
	const query = `
		SELECT
			id,
			board,
			player_x,
			player_o,
			state,
			turn_player,
			winner,
			vs_computer,
			created_at
		FROM games
		WHERE
			(player_x = $1 OR player_o = $1)
			AND state IN ('win', 'draw')
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query finished games: %w", err)
	}
	defer rows.Close()

	var games []model.Game

	for rows.Next() {
		game, err := scanGame(rows)
		if err != nil {
			return nil, fmt.Errorf("scan finished game: %w", err)
		}

		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate finished games: %w", err)
	}

	return games, nil
}

func (r *gameRepository) GetTopPlayers(ctx context.Context, limit int) ([]model.LeaderboardEntry, error) {
	const query = `
		SELECT
			u.id,
			u.login,
			COALESCE(
				SUM(CASE WHEN g.winner = u.id THEN 1 ELSE 0 END)::float
				/
				NULLIF(COUNT(*), 0),
				0
			) AS win_ratio
		FROM users u
		JOIN games g
			ON g.player_x = u.id
			OR g.player_o = u.id
		WHERE g.state IN ('win', 'draw')
		GROUP BY u.id, u.login
		ORDER BY win_ratio DESC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("query top players: %w", err)
	}
	defer rows.Close()

	var players []model.LeaderboardEntry

	for rows.Next() {
		var player model.LeaderboardEntry

		if err := rows.Scan(
			&player.UserID,
			&player.Login,
			&player.WinRatio,
		); err != nil {
			return nil, fmt.Errorf("scan top player: %w", err)
		}

		players = append(players, player)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate top players: %w", err)
	}

	return players, nil
}
