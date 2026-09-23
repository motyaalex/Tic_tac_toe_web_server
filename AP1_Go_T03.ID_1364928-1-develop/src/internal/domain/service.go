package domain

// Интерфейс, который определяет правила игры
type GameService interface {
	GetNextMove(board Board) (row, col int, err error) // Куда ходить (минимакс)
	ValidateBoard(current, updated Board) error        // Корректность хода
	IsGameOver(board Board) bool                       // Есть ли победитель или ничья
}
