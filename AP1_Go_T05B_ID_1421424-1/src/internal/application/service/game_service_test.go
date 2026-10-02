package service

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"tic-tac-toe/internal/domain/model"
	apperrors "tic-tac-toe/internal/errors"
)

// ---------- мок репозитория ----------

type inMemoryGameRepo struct {
	games map[uuid.UUID]model.Game
}

func newInMemoryGameRepo() *inMemoryGameRepo {
	return &inMemoryGameRepo{games: map[uuid.UUID]model.Game{}}
}

func (r *inMemoryGameRepo) Save(_ context.Context, g model.Game) error {
	r.games[g.ID] = g
	return nil
}

func (r *inMemoryGameRepo) Get(_ context.Context, id uuid.UUID) (model.Game, error) {
	g, ok := r.games[id]
	if !ok {
		return model.Game{}, apperrors.ErrNotFound
	}
	return g, nil
}

func (r *inMemoryGameRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.games, id)
	return nil
}

func (r *inMemoryGameRepo) ListWaiting(_ context.Context) ([]model.Game, error) {
	var out []model.Game
	for _, g := range r.games {
		if g.State == model.StateWaitingForPlayers {
			out = append(out, g)
		}
	}
	return out, nil
}

func (r *inMemoryGameRepo) ListFinishedByUser(
	_ context.Context,
	userID uuid.UUID,
) ([]model.Game, error) {
	var out []model.Game

	for _, g := range r.games {
		if g.State != model.StateWin && g.State != model.StateDraw {
			continue
		}

		if (g.PlayerX != nil && *g.PlayerX == userID) ||
			(g.PlayerO != nil && *g.PlayerO == userID) {
			out = append(out, g)
		}
	}

	return out, nil
}

func (r *inMemoryGameRepo) GetTopPlayers(
	_ context.Context,
	limit int,
) ([]model.LeaderboardEntry, error) {
	return []model.LeaderboardEntry{}, nil
}

// ---------- вспомогательные ----------

func newTestGameService() (*gameService, *inMemoryGameRepo) {
	repo := newInMemoryGameRepo()
	return &gameService{repo: repo}, repo
}

// ---------- CreateGame ----------

func TestCreateGameVsComputer(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar := uuid.New()

	g, err := svc.CreateGame(context.Background(), valmerar, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.State != model.StatePlayerTurn {
		t.Errorf("state = %q, want %q", g.State, model.StatePlayerTurn)
	}
	if g.PlayerX == nil || *g.PlayerX != valmerar {
		t.Errorf("playerX = %v, want %v", g.PlayerX, valmerar)
	}
	if g.TurnPlayer == nil || *g.TurnPlayer != valmerar {
		t.Errorf("turnPlayer = %v, want %v", g.TurnPlayer, valmerar)
	}
	if !g.VsComputer {
		t.Error("vsComputer should be true")
	}
}

func TestCreateGameVsPlayer(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar := uuid.New()

	g, err := svc.CreateGame(context.Background(), valmerar, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.State != model.StateWaitingForPlayers {
		t.Errorf("state = %q, want %q", g.State, model.StateWaitingForPlayers)
	}
	if g.PlayerO != nil {
		t.Errorf("playerO should be nil, got %v", g.PlayerO)
	}
	if g.TurnPlayer != nil {
		t.Errorf("turnPlayer should be nil while waiting, got %v", g.TurnPlayer)
	}
}

// ---------- JoinGame ----------

func TestJoinGameSuccess(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar, vasya := uuid.New(), uuid.New()

	g, _ := svc.CreateGame(context.Background(), valmerar, false)
	g, err := svc.JoinGame(context.Background(), g.ID, vasya)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.PlayerO == nil || *g.PlayerO != vasya {
		t.Errorf("playerO = %v, want %v", g.PlayerO, vasya)
	}
	if g.State != model.StatePlayerTurn {
		t.Errorf("state = %q, want %q", g.State, model.StatePlayerTurn)
	}
	if g.TurnPlayer == nil || *g.TurnPlayer != valmerar {
		t.Errorf("turnPlayer should be X (valmerar), got %v", g.TurnPlayer)
	}
}

func TestJoinGameCannotJoinOwn(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar := uuid.New()

	g, _ := svc.CreateGame(context.Background(), valmerar, false)
	_, err := svc.JoinGame(context.Background(), g.ID, valmerar)
	if err == nil {
		t.Fatal("expected error when joining own game")
	}
}

func TestJoinGameNotWaiting(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar, vasya := uuid.New(), uuid.New()

	g, _ := svc.CreateGame(context.Background(), valmerar, false)
	_, _ = svc.JoinGame(context.Background(), g.ID, vasya)

	// третья попытка присоединиться — игра уже не в waiting
	third := uuid.New()
	_, err := svc.JoinGame(context.Background(), g.ID, third)
	if err == nil {
		t.Fatal("expected error when joining a non-waiting game")
	}
}

// ---------- MakeMove ----------

func TestMakeMoveNotYourTurn(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar, vasya := uuid.New(), uuid.New()

	g, _ := svc.CreateGame(context.Background(), valmerar, false)
	g, _ = svc.JoinGame(context.Background(), g.ID, vasya)

	// сейчас ход valmerar, vasya пытается ходить
	_, err := svc.MakeMove(context.Background(), g.ID, vasya, 0, 0)
	if err == nil {
		t.Fatal("expected error: not your turn")
	}
}

func TestMakeMoveOutOfBounds(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar := uuid.New()

	g, _ := svc.CreateGame(context.Background(), valmerar, true)
	for _, tc := range []struct{ r, c int }{
		{-1, 0}, {0, -1}, {3, 0}, {0, 3}, {99, 99},
	} {
		_, err := svc.MakeMove(context.Background(), g.ID, valmerar, tc.r, tc.c)
		if err == nil {
			t.Errorf("expected error for (%d,%d)", tc.r, tc.c)
		}
	}
}

func TestMakeMoveCellOccupied(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar, vasya := uuid.New(), uuid.New()

	g, _ := svc.CreateGame(context.Background(), valmerar, false)
	g, _ = svc.JoinGame(context.Background(), g.ID, vasya)

	// valmerar (X) ходит (0,0)
	g, err := svc.MakeMove(context.Background(), g.ID, valmerar, 0, 0)
	if err != nil {
		t.Fatalf("first move: %v", err)
	}

	// vasya (O) пытается сходить в ту же клетку
	_, err = svc.MakeMove(context.Background(), g.ID, vasya, 0, 0)
	if err == nil {
		t.Fatal("expected error: cell occupied")
	}
}

func TestMakeMoveByStranger(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar, vasya, stranger := uuid.New(), uuid.New(), uuid.New()

	g, _ := svc.CreateGame(context.Background(), valmerar, false)
	g, _ = svc.JoinGame(context.Background(), g.ID, vasya)

	_, err := svc.MakeMove(context.Background(), g.ID, stranger, 0, 0)
	if err == nil {
		t.Fatal("expected error: stranger cannot play")
	}
}

// ---------- полный сценарий: победа X ----------

func TestFullGameXWins(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar, vasya := uuid.New(), uuid.New()

	g, _ := svc.CreateGame(context.Background(), valmerar, false)
	g, _ = svc.JoinGame(context.Background(), g.ID, vasya)

	// X: (0,0) (0,1) (0,2)   O: (1,0) (1,1)
	moves := []struct {
		user uuid.UUID
		r, c int
	}{
		{valmerar, 0, 0}, {vasya, 1, 0},
		{valmerar, 0, 1}, {vasya, 1, 1},
		{valmerar, 0, 2}, // победный
	}
	for i, m := range moves {
		var err error
		g, err = svc.MakeMove(context.Background(), g.ID, m.user, m.r, m.c)
		if err != nil {
			t.Fatalf("move %d: %v", i, err)
		}
	}

	if g.State != model.StateWin {
		t.Errorf("state = %q, want %q", g.State, model.StateWin)
	}
	if g.Winner == nil || *g.Winner != valmerar {
		t.Errorf("winner = %v, want %v", g.Winner, valmerar)
	}
	if g.TurnPlayer != nil {
		t.Errorf("turnPlayer should be nil on finish, got %v", g.TurnPlayer)
	}
}

// ---------- полный сценарий: ничья ----------

func TestFullGameDraw(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar, vasya := uuid.New(), uuid.New()

	g, _ := svc.CreateGame(context.Background(), valmerar, false)
	g, _ = svc.JoinGame(context.Background(), g.ID, vasya)

	// X: (0,0) (0,2) (1,0) (2,1) (2,2)
	// O: (0,1) (1,1) (1,2) (2,0)
	moves := []struct {
		user uuid.UUID
		r, c int
	}{
		{valmerar, 0, 0}, {vasya, 0, 1},
		{valmerar, 0, 2}, {vasya, 1, 1},
		{valmerar, 1, 0}, {vasya, 1, 2},
		{valmerar, 2, 1}, {vasya, 2, 0},
		{valmerar, 2, 2},
	}
	for i, m := range moves {
		var err error
		g, err = svc.MakeMove(context.Background(), g.ID, m.user, m.r, m.c)
		if err != nil {
			t.Fatalf("move %d: %v", i, err)
		}
	}

	if g.State != model.StateDraw {
		t.Errorf("state = %q, want %q", g.State, model.StateDraw)
	}
	if g.Winner != nil {
		t.Errorf("winner should be nil on draw, got %v", g.Winner)
	}
}

// ---------- игра с компьютером ----------

func TestComputerMovesAfterPlayer(t *testing.T) {
	svc, _ := newTestGameService()
	valmerar := uuid.New()

	g, _ := svc.CreateGame(context.Background(), valmerar, true)
	g, err := svc.MakeMove(context.Background(), g.ID, valmerar, 1, 1)
	if err != nil {
		t.Fatalf("move: %v", err)
	}
	// после хода человека компьютер должен был сходить сам
	computerCells := 0
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if g.Board[r][c] == 2 {
				computerCells++
			}
		}
	}
	if computerCells != 1 {
		t.Errorf("computer cells = %d, want 1", computerCells)
	}
	// теперь снова ход человека
	if g.State != model.StatePlayerTurn {
		t.Errorf("state = %q, want %q", g.State, model.StatePlayerTurn)
	}
	if g.TurnPlayer == nil || *g.TurnPlayer != valmerar {
		t.Errorf("turnPlayer should be valmerar again, got %v", g.TurnPlayer)
	}
}

// ---------- юнит-тесты на чистые функции ----------

func TestWinnerAllLines(t *testing.T) {
	tests := []struct {
		name  string
		board model.Board
		want  int
	}{
		{"empty", model.Board{}, 0},
		{"row0 X", model.Board{{1, 1, 1}, {0, 0, 0}, {0, 0, 0}}, 1},
		{"row1 X", model.Board{{0, 0, 0}, {1, 1, 1}, {0, 0, 0}}, 1},
		{"row2 O", model.Board{{0, 0, 0}, {0, 0, 0}, {2, 2, 2}}, 2},
		{"col0 X", model.Board{{1, 0, 0}, {1, 0, 0}, {1, 0, 0}}, 1},
		{"col2 O", model.Board{{0, 0, 2}, {0, 0, 2}, {0, 0, 2}}, 2},
		{"diag X", model.Board{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}, 1},
		{"antidiag O", model.Board{{0, 0, 2}, {0, 2, 0}, {2, 0, 0}}, 2},
		{"no winner", model.Board{{1, 2, 1}, {1, 2, 2}, {2, 1, 1}}, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := winner(tc.board); got != tc.want {
				t.Errorf("winner() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestBoardFull(t *testing.T) {
	if boardFull(model.Board{}) {
		t.Error("empty board should not be full")
	}
	full := model.Board{
		{1, 2, 1},
		{2, 1, 2},
		{2, 1, 1},
	}
	if !boardFull(full) {
		t.Error("full board should be full")
	}
}

func TestSymbolFor(t *testing.T) {
	valmerar, vasya, stranger := uuid.New(), uuid.New(), uuid.New()
	g := model.Game{PlayerX: &valmerar, PlayerO: &vasya}

	if s := g.SymbolFor(valmerar); s != model.SymbolX {
		t.Errorf("valmerar: got %q, want X", s)
	}
	if s := g.SymbolFor(vasya); s != model.SymbolO {
		t.Errorf("vasya: got %q, want O", s)
	}
	if s := g.SymbolFor(stranger); s != "" {
		t.Errorf("stranger: got %q, want empty", s)
	}
}
