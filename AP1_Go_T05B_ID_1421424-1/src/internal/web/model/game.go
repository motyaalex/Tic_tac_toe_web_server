package model

import "github.com/google/uuid"

type GameState string

const (
	StateWaitingForPlayers GameState = "waiting_for_players"
	StatePlayerTurn        GameState = "player_turn"
	StateDraw              GameState = "draw"
	StateWin               GameState = "win"
)

type Symbol string

const (
	SymbolX Symbol = "X"
	SymbolO Symbol = "O"
)

var ComputerID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type Game struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	Board      Board      `json:"board" db:"board"`
	PlayerX    *uuid.UUID `json:"playerX,omitempty" db:"player_x"`
	PlayerO    *uuid.UUID `json:"playerO,omitempty" db:"player_o"`
	State      GameState  `json:"state" db:"state"`
	TurnPlayer *uuid.UUID `json:"turnPlayer,omitempty" db:"turn_player"`
	Winner     *uuid.UUID `json:"winner,omitempty" db:"winner"`
	VsComputer bool       `json:"vsComputer" db:"vs_computer"`
}

type Board [3][3]int

func (b Board) IsEmpty(r, c int) bool { return b[r][c] == 0 }

func NewGame(playerX uuid.UUID, vsComputer bool) Game {
	g := Game{
		ID:         uuid.New(),
		PlayerX:    &playerX,
		State:      StateWaitingForPlayers,
		VsComputer: vsComputer,
	}
	if vsComputer {
		g.PlayerO = &ComputerID
		g.State = StatePlayerTurn
		g.TurnPlayer = &playerX
	}
	return g
}

func (g Game) SymbolFor(userID uuid.UUID) Symbol {
	if g.PlayerX != nil && *g.PlayerX == userID {
		return SymbolX
	}
	if g.PlayerO != nil && *g.PlayerO == userID {
		return SymbolO
	}
	return ""
}
