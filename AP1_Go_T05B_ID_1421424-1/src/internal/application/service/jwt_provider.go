package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"tic-tac-toe/internal/config"
	"tic-tac-toe/internal/domain/model"
)

type JwtProvider struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type jwtClaims struct {
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}

func NewJwtProvider(cfg config.Config) *JwtProvider {
	return &JwtProvider{
		secret:          []byte(cfg.JWTSecret),
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}
}

func (p *JwtProvider) GenerateAccessToken(user model.User) (string, error) {
	return p.generateToken(user.ID, "access", p.accessTokenTTL)
}

func (p *JwtProvider) GenerateRefreshToken(user model.User) (string, error) {
	return p.generateToken(user.ID, "refresh", p.refreshTokenTTL)
}

func (p *JwtProvider) generateToken(
	userID uuid.UUID,
	tokenType string,
	ttl time.Duration,
) (string, error) {
	now := time.Now()

	claims := jwtClaims{
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(p.secret)
}

func (p *JwtProvider) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	return p.validateToken(tokenString, "access")
}

func (p *JwtProvider) ValidateRefreshToken(tokenString string) (uuid.UUID, error) {
	return p.validateToken(tokenString, "refresh")
}

func (p *JwtProvider) validateToken(
	tokenString string,
	expectedType string,
) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return p.secret, nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}

	if claims.TokenType != expectedType {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
