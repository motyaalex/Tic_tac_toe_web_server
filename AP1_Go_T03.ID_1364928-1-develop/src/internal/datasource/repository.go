package datasource

import (
	"tic-tac-toe/internal/domain"
)

// Структура, которая РЕАЛИЗУЕТ интерфейс domain.GameRepository
type gameRepository struct {
	storage *GameStorage // конкретное хранилище
}

func NewGameRepository(storage *GameStorage) domain.GameRepository {
	return &gameRepository{storage: storage}
}

// Реализация метода Save, сохраняет текущую игру с ID и состоянием поля
func (r *gameRepository) Save(game *domain.Game) error {
	r.storage.games.Store(game.ID, ToDTO(game))
	return nil
}

// Реализация метода Get, находит игру с ID и состоянием поля
func (r *gameRepository) Get(id string) (*domain.Game, error) {
	val, ok := r.storage.games.Load(id)
	if !ok {
		return nil, domain.ErrGameNotFound
	}
	return ToDomain(val.(*GameDTO)), nil
}
