package model

import "github.com/google/uuid"

type LeaderboardEntry struct {
	UserID   uuid.UUID `json:"userId"`
	Login    string    `json:"login"`
	WinRatio float64   `json:"winRatio"`
}
