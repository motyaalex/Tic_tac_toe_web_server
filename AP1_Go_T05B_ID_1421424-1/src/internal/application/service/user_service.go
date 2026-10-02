package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"tic-tac-toe/internal/datasource/repository"
	"tic-tac-toe/internal/domain/model"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req model.SignUpRequest) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	u := model.User{
		ID:           uuid.New(),
		Login:        req.Login,
		PasswordHash: string(hash),
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (s *UserService) FindByLogin(ctx context.Context, login string) (model.User, error) {
	return s.repo.GetByLogin(ctx, login)
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	return s.repo.GetByID(ctx, id)
}
