package domain

import "errors"

// Структура ошибок с текстом
var (
	ErrGameNotFound    = errors.New("game not found")
	ErrInvalidMove     = errors.New("invalid move")
	ErrBoardChanged    = errors.New("board was changed illegally")
	ErrGameAlreadyOver = errors.New("game is already over")
)
