// internal/domain/model/board.go
package model

import "errors"

// Board представляет игровое поле 3x3
// 0 = пусто, 1 = игрок (X), 2 = компьютер (O)
type Board [3][3]int

// NewBoard создаёт пустое поле
func NewBoard() Board {
	return Board{}
}

// MakeMove делает ход (row, col, player)
// Возвращает ошибку, если ход некорректен
func (b *Board) MakeMove(row, col, player int) error {
	if row < 0 || row > 2 || col < 0 || col > 2 {
		return errors.New("некорректные координаты")
	}
	if b[row][col] != 0 {
		return errors.New("клетка занята")
	}
	b[row][col] = player
	return nil
}
func (b Board) IsEmpty(r, c int) bool {
	return b[r][c] == 0
}

// IsFull проверяет, заполнено ли поле
func (b Board) IsFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

// GetWinner возвращает победителя (1 — игрок X, 2 — компьютер O) или 0, если победителя нет
func (b Board) GetWinner() int {
	// Проверка строк (горизонтали)
	for row := 0; row < 3; row++ {
		if b[row][0] != 0 && b[row][0] == b[row][1] && b[row][1] == b[row][2] {
			return b[row][0]
		}
	}

	// Проверка столбцов (вертикали)
	for col := 0; col < 3; col++ {
		if b[0][col] != 0 && b[0][col] == b[1][col] && b[1][col] == b[2][col] {
			return b[0][col]
		}
	}

	// Проверка главной диагонали (слева сверху → вправо вниз)
	if b[0][0] != 0 && b[0][0] == b[1][1] && b[1][1] == b[2][2] {
		return b[0][0]
	}

	// Проверка побочной диагонали (справа сверху → влево вниз)
	if b[0][2] != 0 && b[0][2] == b[1][1] && b[1][1] == b[2][0] {
		return b[0][2]
	}

	return 0 // нет победителя
}
