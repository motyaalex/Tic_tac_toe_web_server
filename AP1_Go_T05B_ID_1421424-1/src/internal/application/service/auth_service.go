package service

import (
	"context"
	"fmt"

	"tic-tac-toe/internal/domain/model"
	apperrors "tic-tac-toe/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct { // Переделано через jwt
	users *UserService
	jwt   *JwtProvider
}

func NewAuthService(users *UserService, jwt *JwtProvider) *AuthService { // Переделано через jwt
	return &AuthService{
		users: users,
		jwt:   jwt,
	}
}

// SignUp регистрирует нового пользователя.
func (s *AuthService) SignUp(ctx context.Context, req model.SignUpRequest) error {
	if err := ValidateCredentials(req.Login, req.Password); err != nil {
		return err
	}
	_, err := s.users.CreateUser(ctx, req)
	return err
}

// SignIn проверяет логин/пароль и возвращает UUID пользователя. ПРЕДЕЛАНО! Новый метод, возвращает jwt токены
func (s *AuthService) SignIn(ctx context.Context, req model.JwtRequest) (model.JwtResponse, error) {
	user, err := s.users.FindByLogin(ctx, req.Login)
	if err != nil {
		return model.JwtResponse{}, apperrors.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return model.JwtResponse{}, apperrors.ErrUnauthorized
	}

	accessToken, err := s.jwt.GenerateAccessToken(user)
	if err != nil {
		return model.JwtResponse{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user)
	if err != nil {
		return model.JwtResponse{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return model.JwtResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Токены для jwt, создаёт новый на основании refresh token
func (s *AuthService) RefreshAccessToken(ctx context.Context, req model.RefreshJwtRequest) (model.JwtResponse, error) {
	userID, err := s.jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return model.JwtResponse{}, apperrors.ErrUnauthorized
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return model.JwtResponse{}, apperrors.ErrUnauthorized
	}

	accessToken, err := s.jwt.GenerateAccessToken(user)
	if err != nil {
		return model.JwtResponse{}, fmt.Errorf("generate access token: %w", err)
	}

	return model.JwtResponse{
		Type:        "Bearer",
		AccessToken: accessToken,
	}, nil
}

// Токены для jwt, создаёт новый на основании существующего refresh token
func (s *AuthService) RefreshRefreshToken(ctx context.Context, req model.RefreshJwtRequest) (model.JwtResponse, error) {
	userID, err := s.jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return model.JwtResponse{}, apperrors.ErrUnauthorized
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return model.JwtResponse{}, apperrors.ErrUnauthorized
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user)
	if err != nil {
		return model.JwtResponse{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return model.JwtResponse{
		Type:         "Bearer",
		RefreshToken: refreshToken,
	}, nil
}

// ValidateCredentials проверяет логин и пароль до выполнения защищённого запроса.
func ValidateCredentials(login, password string) error {
	if len(login) < 3 {
		return fmt.Errorf("login too short")
	}
	if len(password) < 5 {
		return fmt.Errorf("password too short")
	}
	return nil
}
