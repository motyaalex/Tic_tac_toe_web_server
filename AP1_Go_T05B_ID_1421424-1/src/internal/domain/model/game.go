package model

import (
	"time"

	"github.com/google/uuid"
)

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

// internal/domain/model/game.go
var ComputerID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type Game struct {
	ID         uuid.UUID  `json:"id"                     db:"id"`
	Board      Board      `json:"board"                  db:"board"`
	PlayerX    *uuid.UUID `json:"playerX,omitempty"      db:"player_x"`
	PlayerO    *uuid.UUID `json:"playerO,omitempty"      db:"player_o"`
	State      GameState  `json:"state"                  db:"state"`
	TurnPlayer *uuid.UUID `json:"turnPlayer,omitempty"   db:"turn_player"`
	Winner     *uuid.UUID `json:"winner,omitempty"       db:"winner"`
	VsComputer bool       `json:"vsComputer"             db:"vs_computer"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
}

func NewGame(playerX uuid.UUID, vsComputer bool) Game {
	g := Game{
		ID:         uuid.New(),
		PlayerX:    &playerX,
		State:      StateWaitingForPlayers,
		VsComputer: vsComputer,
		CreatedAt:  time.Now(),
	}
	if vsComputer {
		// игрок сразу играет против компьютера, второй игрок не нужен
		g.PlayerO = &ComputerID
		g.State = StatePlayerTurn
		g.TurnPlayer = &playerX
	}
	return g
}

// SymbolFor возвращает значок игрока (или "", если игрок не в этой игре).
func (g Game) SymbolFor(userID uuid.UUID) Symbol {
	if g.PlayerX != nil && *g.PlayerX == userID {
		return SymbolX
	}
	if g.PlayerO != nil && *g.PlayerO == userID {
		return SymbolO
	}
	return ""
}
