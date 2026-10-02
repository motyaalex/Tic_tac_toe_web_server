package model

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `json:"id"    db:"id"`
	Login        string    `json:"login" db:"login"`
	PasswordHash string    `json:"-"     db:"password_hash"` // "-" чтобы не отдавать в JSON
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JwtRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JwtResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshJwtRequest struct {
	RefreshToken string `json:"refreshToken"`
}
