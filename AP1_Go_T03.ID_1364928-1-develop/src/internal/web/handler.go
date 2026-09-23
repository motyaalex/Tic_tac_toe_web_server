package web

import (
	"encoding/json"
	"net/http"
	"strings"
	"tic-tac-toe/internal/domain"
)

type Handler struct {
	service domain.GameService
	repo    domain.GameRepository // Использует интерфейс domain.GameRepository
}

func NewHandler(service domain.GameService, repo domain.GameRepository) *Handler {
	return &Handler{
		service: service,
		repo:    repo,
	}
}

// НОВЫЙ МЕТОД: создание новой игры
func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Создаём новую игру (UUID генерируется в NewGame())
	game := domain.NewGame()

	// Сохраняем в хранилище
	if err := h.repo.Save(game); err != nil {
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ с ID новой игры
	response := CreateGameResponse{
		ID:      game.ID,
		Message: "Game created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DTO для ответа при создании игры
type CreateGameResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID игры из URL
	path := strings.TrimPrefix(r.URL.Path, "/game/")
	if path == "" {
		http.Error(w, "Game ID required", http.StatusBadRequest)
		return
	}
	gameID := strings.Split(path, "/")[0]

	// Декодируем запрос
	var req GameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем текущую игру
	currentGame, err := h.repo.Get(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Проверяем, не закончена ли игра
	if h.service.IsGameOver(currentGame.Board) {
		resp := FromDomain(currentGame, true, "Game is already over")
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Валидируем ход игрока
	updatedGame := ToDomain(gameID, req)
	if err := h.service.ValidateBoard(currentGame.Board, updatedGame.Board); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, не выиграл ли игрок после своего хода
	if h.service.IsGameOver(updatedGame.Board) {
		h.repo.Save(updatedGame)
		resp := FromDomain(updatedGame, true, "Player wins!")
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Ход компьютера
	row, col, err := h.service.GetNextMove(updatedGame.Board)
	if err == nil {
		updatedGame.Board[row][col] = 2 // Ход компьютера
	}

	// Проверяем окончание игры
	gameOver := h.service.IsGameOver(updatedGame.Board)
	message := ""
	if gameOver {
		if err != nil {
			message = "Draw!"
		} else {
			message = "Computer wins!"
		}
	}

	// Сохраняем обновленную игру
	h.repo.Save(updatedGame)

	// Отправляем ответ
	resp := FromDomain(updatedGame, gameOver, message)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
