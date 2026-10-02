package service

import (
	"context"

	"github.com/google/uuid"

	"tic-tac-toe/internal/datasource/repository"
	"tic-tac-toe/internal/domain/model"
	domainService "tic-tac-toe/internal/domain/service"
	apperrors "tic-tac-toe/internal/errors"
)

type gameService struct {
	repo repository.GameRepository
}

func NewGameService(repo repository.GameRepository) domainService.GameService {
	return &gameService{repo: repo}
}

func (s *gameService) CreateGame(ctx context.Context, userID uuid.UUID, vsComputer bool) (model.Game, error) {
	g := model.NewGame(userID, vsComputer)
	if err := s.repo.Save(ctx, g); err != nil {
		return model.Game{}, err
	}
	return g, nil
}

func (s *gameService) JoinGame(ctx context.Context, gameID, userID uuid.UUID) (model.Game, error) {
	g, err := s.repo.Get(ctx, gameID)
	if err != nil {
		return model.Game{}, err
	}
	if g.State != model.StateWaitingForPlayers {
		return model.Game{}, apperrors.ErrValidation
	}
	if g.PlayerX != nil && *g.PlayerX == userID {
		return model.Game{}, apperrors.ErrValidation
	}
	g.PlayerO = &userID
	g.State = model.StatePlayerTurn
	g.TurnPlayer = g.PlayerX
	if err := s.repo.Save(ctx, g); err != nil {
		return model.Game{}, err
	}
	return g, nil
}

func (s *gameService) MakeMove(ctx context.Context, gameID, userID uuid.UUID, row, col int) (model.Game, error) {
	g, err := s.repo.Get(ctx, gameID)
	if err != nil {
		return model.Game{}, err
	}
	if g.State != model.StatePlayerTurn {
		return model.Game{}, apperrors.ErrValidation
	}
	if g.TurnPlayer == nil || *g.TurnPlayer != userID {
		return model.Game{}, apperrors.ErrValidation
	}
	if row < 0 || row > 2 || col < 0 || col > 2 {
		return model.Game{}, apperrors.ErrValidation
	}
	if !g.Board.IsEmpty(row, col) {
		return model.Game{}, apperrors.ErrValidation
	}

	symbol := g.SymbolFor(userID)
	if symbol == "" {
		return model.Game{}, apperrors.ErrValidation
	}

	applyMove(&g, row, col, symbol)
	updateState(&g)

	// Если игра против компьютера и сейчас его ход — компьютер ходит сам.
	if g.VsComputer && g.State == model.StatePlayerTurn {
		s.computerMove(&g)
		updateState(&g)
	}

	if err := s.repo.Save(ctx, g); err != nil {
		return model.Game{}, err
	}
	return g, nil
}

func (s *gameService) GetGame(ctx context.Context, gameID uuid.UUID) (model.Game, error) {
	return s.repo.Get(ctx, gameID)
}

func (s *gameService) ListWaiting(ctx context.Context) ([]model.Game, error) {
	return s.repo.ListWaiting(ctx)
}

func (s *gameService) ListFinishedByUser(ctx context.Context, userID uuid.UUID) ([]model.Game, error) {
	return s.repo.ListFinishedByUser(ctx, userID)
}

func (s *gameService) DeleteGame(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.Get(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

// --- helpers ---

func applyMove(g *model.Game, r, c int, sym model.Symbol) {
	val := 1
	if sym == model.SymbolO {
		val = 2
	}
	g.Board[r][c] = val
}

func updateState(g *model.Game) {
	if w := winner(g.Board); w != 0 {
		g.State = model.StateWin
		g.TurnPlayer = nil
		if w == 1 {
			g.Winner = g.PlayerX
		} else {
			g.Winner = g.PlayerO
		}
		return
	}
	if boardFull(g.Board) {
		g.State = model.StateDraw
		g.TurnPlayer = nil
		return
	}
	// ход соперника
	g.State = model.StatePlayerTurn
	if g.TurnPlayer != nil && g.PlayerX != nil && *g.TurnPlayer == *g.PlayerX {
		g.TurnPlayer = g.PlayerO
	} else {
		g.TurnPlayer = g.PlayerX
	}
}

func boardFull(b model.Board) bool {
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if b[r][c] == 0 {
				return false
			}
		}
	}
	return true
}

func winner(b model.Board) int {
	winLines := [8][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}}, {{1, 0}, {1, 1}, {1, 2}}, {{2, 0}, {2, 1}, {2, 2}},
		{{0, 0}, {1, 0}, {2, 0}}, {{0, 1}, {1, 1}, {2, 1}}, {{0, 2}, {1, 2}, {2, 2}},
		{{0, 0}, {1, 1}, {2, 2}}, {{0, 2}, {1, 1}, {2, 0}},
	}
	for _, l := range winLines {
		a1, a2, a3 := b[l[0][0]][l[0][1]], b[l[1][0]][l[1][1]], b[l[2][0]][l[2][1]]
		if a1 != 0 && a1 == a2 && a2 == a3 {
			return a1
		}
	}
	return 0
}

// computerMove — компьютер ходит O, используя минимакс.
func (s *gameService) computerMove(g *model.Game) {
	board := g.Board // копия, чтобы minimax не портил оригинал
	r, c := bestMove(board)
	if r == -1 {
		return
	}
	g.Board[r][c] = 2
}

func bestMove(board model.Board) (bestRow, bestCol int) {
	bestScore := -1 << 30
	bestRow, bestCol = -1, -1
	for r := range board {
		for c := range board[r] {
			if board[r][c] != 0 {
				continue
			}
			board[r][c] = 2
			score := minimax(board, false, -1<<30, 1<<30, 0)
			board[r][c] = 0
			if score > bestScore {
				bestScore = score
				bestRow, bestCol = r, c
			}
		}
	}
	return
}

func minimax(board model.Board, isComputerTurn bool, alpha, beta, depth int) int {
	if w := winner(board); w != 0 {
		if w == 2 {
			return 10 - depth
		}
		return -10 + depth
	}
	empty := 0
	for r := range board {
		for c := range board[r] {
			if board[r][c] == 0 {
				empty++
			}
		}
	}
	if empty == 0 {
		return 0
	}

	if isComputerTurn {
		best := -1 << 30
		for r := range board {
			for c := range board[r] {
				if board[r][c] != 0 {
					continue
				}
				board[r][c] = 2
				best = max(best, minimax(board, false, alpha, beta, depth+1))
				board[r][c] = 0
				alpha = max(alpha, best)
				if beta <= alpha {
					return best
				}
			}
		}
		return best
	}

	best := 1 << 30
	for r := range board {
		for c := range board[r] {
			if board[r][c] != 0 {
				continue
			}
			board[r][c] = 1
			best = min(best, minimax(board, true, alpha, beta, depth+1))
			board[r][c] = 0
			beta = min(beta, best)
			if beta <= alpha {
				return best
			}
		}
	}
	return best
}

func (s *gameService) GetTopPlayers(ctx context.Context, limit int) ([]model.LeaderboardEntry, error) {
	return s.repo.GetTopPlayers(ctx, limit)
}
