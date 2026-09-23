package domain

// Интерфейс, который сохраняет и находит игру по UUID
type GameRepository interface {
	Save(game *Game) error
	Get(id string) (*Game, error)
}
