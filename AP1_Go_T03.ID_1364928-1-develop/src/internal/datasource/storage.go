package datasource

import (
	"sync"
)

// Класс-хранилище для текущих игр
type GameStorage struct {
	games sync.Map // Потокобезопасные коллекции хранилищ игр
}

func NewGameStorage() *GameStorage {
	return &GameStorage{}
}

type GameDTO struct {
	ID    string
	Board [3][3]int // Модель игрового поля
}
