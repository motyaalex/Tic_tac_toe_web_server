package web

// для отправки запроса от клиента в Domain
type GameRequest struct {
	Board [3][3]int `json:"board"` // Модель игрового поля
}

// для отправки ответа клиенту после обработки хода
type GameResponse struct {
	ID      string    `json:"id"`
	Board   [3][3]int `json:"board"` // Модель игрового поля
	Over    bool      `json:"over"`
	Message string    `json:"message,omitempty"`
}
