package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	domainService "tic-tac-toe/internal/domain/service"
	"tic-tac-toe/internal/middleware"
)

type GameHandler struct {
	svc domainService.GameService
}

func NewGameHandler(svc domainService.GameService) *GameHandler {
	return &GameHandler{svc: svc}
}

func (h *GameHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /game", h.Create)
	mux.HandleFunc("GET /games", h.List)
	mux.HandleFunc("GET /game/{id}", h.Get)
	mux.HandleFunc("POST /game/{id}/join", h.Join)
	mux.HandleFunc("POST /game/{id}/move", h.Move)
	mux.HandleFunc("GET /history", h.History)
	mux.HandleFunc("GET /leaderboard", h.Leaderboard)
}

type moveRequest struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func (h *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	vsComputer := r.URL.Query().Get("vsComputer") == "true"

	g, err := h.svc.CreateGame(r.Context(), userID, vsComputer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, g)
}

func (h *GameHandler) List(w http.ResponseWriter, r *http.Request) {
	games, err := h.svc.ListWaiting(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, games)
}

func (h *GameHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	g, err := h.svc.GetGame(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func (h *GameHandler) Join(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	g, err := h.svc.JoinGame(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func (h *GameHandler) Move(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req moveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	g, err := h.svc.MakeMove(r.Context(), id, userID, req.Row, req.Col)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func (h *GameHandler) History(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	games, err := h.svc.ListFinishedByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, games)
}

func (h *GameHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	nStr := r.URL.Query().Get("n")

	if nStr == "" {
		http.Error(w, "missing n", http.StatusBadRequest)
		return
	}

	n, err := strconv.Atoi(nStr)
	if err != nil || n <= 0 {
		http.Error(w, "invalid n", http.StatusBadRequest)
		return
	}

	players, err := h.svc.GetTopPlayers(r.Context(), n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, players)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
