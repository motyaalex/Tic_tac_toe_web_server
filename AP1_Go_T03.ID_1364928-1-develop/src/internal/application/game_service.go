package application

import (
	"tic-tac-toe/internal/domain"
)

// Использует интерфейс domain.GameRepository
type gameService struct {
	repo domain.GameRepository // Принимает репозиторий
}

func NewGameService(repo domain.GameRepository) domain.GameService {
	return &gameService{repo: repo}
}

// Минимакс алгоритм (упрощенная версия) - Minimax
func (s *gameService) GetNextMove(board domain.Board) (int, int, error) {
	// Проверяем, не закончена ли игра
	if s.IsGameOver(board) {
		return -1, -1, domain.ErrGameAlreadyOver
	}

	// Находим первый пустой слот для хода компьютера
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				return i, j, nil
			}
		}
	}
	return -1, -1, domain.ErrInvalidMove
}

func (s *gameService) ValidateBoard(current, updated domain.Board) error {
	// Проверяем, что изменился только один слот и он стал 1 (ход игрока)
	changes := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if current[i][j] != updated[i][j] {
				changes++
				if updated[i][j] != 1 || current[i][j] != 0 {
					return domain.ErrInvalidMove
				}
			}
		}
	}
	if changes != 1 {
		return domain.ErrInvalidMove
	}
	return nil
}

func (s *gameService) IsGameOver(board domain.Board) bool {
	// Проверка горизонталей
	for i := 0; i < 3; i++ {
		if board[i][0] != 0 && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			return true
		}
	}
	// Проверка вертикалей
	for j := 0; j < 3; j++ {
		if board[0][j] != 0 && board[0][j] == board[1][j] && board[1][j] == board[2][j] {
			return true
		}
	}
	// Проверка диагоналей
	if board[0][0] != 0 && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return true
	}
	if board[0][2] != 0 && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return true
	}

	// Проверка на ничью (нет пустых клеток)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				return false
			}
		}
	}
	return true
}
