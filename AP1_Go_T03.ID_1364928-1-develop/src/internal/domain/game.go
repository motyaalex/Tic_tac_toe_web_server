package domain

import "github.com/google/uuid"

//Струтура с ID игры и состоянием поля
type Game struct {
	ID    string
	Board Board
}

//Функция создания новой игры с ID и полем
func NewGame() *Game {
	return &Game{
		ID:    uuid.New().String(),
		Board: Board{},
	}
}
