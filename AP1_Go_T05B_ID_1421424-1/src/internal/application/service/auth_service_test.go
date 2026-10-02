package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"tic-tac-toe/internal/config"
	"tic-tac-toe/internal/domain/model"
	apperrors "tic-tac-toe/internal/errors"
)

// ---------- мок UserRepository ----------

type inMemoryUserRepo struct {
	byLogin map[string]model.User
	byID    map[uuid.UUID]model.User
}

func newInMemoryUserRepo() *inMemoryUserRepo {
	return &inMemoryUserRepo{
		byLogin: map[string]model.User{},
		byID:    map[uuid.UUID]model.User{},
	}
}

func (r *inMemoryUserRepo) Create(_ context.Context, u model.User) error {
	if _, exists := r.byLogin[u.Login]; exists {
		return apperrors.ErrValidation
	}

	r.byLogin[u.Login] = u
	r.byID[u.ID] = u

	return nil
}

func (r *inMemoryUserRepo) GetByLogin(
	_ context.Context,
	login string,
) (model.User, error) {
	u, ok := r.byLogin[login]
	if !ok {
		return model.User{}, apperrors.ErrNotFound
	}

	return u, nil
}

func (r *inMemoryUserRepo) GetByID(
	_ context.Context,
	id uuid.UUID,
) (model.User, error) {
	u, ok := r.byID[id]
	if !ok {
		return model.User{}, apperrors.ErrNotFound
	}

	return u, nil
}

// ---------- вспомогательные ----------

func newTestAuthService() *AuthService {
	repo := newInMemoryUserRepo()
	userSvc := NewUserService(repo)

	jwtProvider := NewJwtProvider(config.Config{
		JWTSecret:       "test-secret",
		AccessTokenTTL:  time.Hour,
		RefreshTokenTTL: 24 * time.Hour,
	})

	return NewAuthService(userSvc, jwtProvider)
}

// ---------- SignUp ----------

func TestSignUpSuccess(t *testing.T) {
	auth := newTestAuthService()

	err := auth.SignUp(context.Background(), model.SignUpRequest{
		Login:    "valmerar",
		Password: "secret123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSignUpShortLogin(t *testing.T) {
	auth := newTestAuthService()

	err := auth.SignUp(context.Background(), model.SignUpRequest{
		Login:    "ab",
		Password: "secret123",
	})

	if err == nil {
		t.Fatal("expected error: login too short")
	}
}

func TestSignUpShortPassword(t *testing.T) {
	auth := newTestAuthService()

	err := auth.SignUp(context.Background(), model.SignUpRequest{
		Login:    "valmerar",
		Password: "12",
	})

	if err == nil {
		t.Fatal("expected error: password too short")
	}
}

func TestSignUpDuplicate(t *testing.T) {
	auth := newTestAuthService()

	req := model.SignUpRequest{
		Login:    "valmerar",
		Password: "secret123",
	}

	if err := auth.SignUp(context.Background(), req); err != nil {
		t.Fatalf("first signup: %v", err)
	}

	if err := auth.SignUp(context.Background(), req); err == nil {
		t.Fatal("expected error: duplicate login")
	}
}

// ---------- SignIn ----------

func TestSignInSuccess(t *testing.T) {
	auth := newTestAuthService()

	_ = auth.SignUp(context.Background(), model.SignUpRequest{
		Login:    "valmerar",
		Password: "secret123",
	})

	response, err := auth.SignIn(
		context.Background(),
		model.JwtRequest{
			Login:    "valmerar",
			Password: "secret123",
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.AccessToken == "" {
		t.Error("access token should not be empty")
	}

	if response.RefreshToken == "" {
		t.Error("refresh token should not be empty")
	}

	if response.Type != "Bearer" {
		t.Errorf("type = %q, want Bearer", response.Type)
	}
}

func TestSignInWrongPassword(t *testing.T) {
	auth := newTestAuthService()

	_ = auth.SignUp(context.Background(), model.SignUpRequest{
		Login:    "valmerar",
		Password: "secret123",
	})

	_, err := auth.SignIn(
		context.Background(),
		model.JwtRequest{
			Login:    "valmerar",
			Password: "WRONG",
		},
	)

	if err == nil {
		t.Fatal("expected error: wrong password")
	}
}

func TestSignInUnknownUser(t *testing.T) {
	auth := newTestAuthService()

	_, err := auth.SignIn(
		context.Background(),
		model.JwtRequest{
			Login:    "nobody",
			Password: "whatever",
		},
	)

	if err == nil {
		t.Fatal("expected error: unknown user")
	}
}

// ---------- пароль хранится хэшем ----------

func TestPasswordIsHashed(t *testing.T) {
	repo := newInMemoryUserRepo()
	userSvc := NewUserService(repo)

	_, err := userSvc.CreateUser(
		context.Background(),
		model.SignUpRequest{
			Login:    "valmerar",
			Password: "secret123",
		},
	)

	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	u, _ := repo.GetByLogin(
		context.Background(),
		"valmerar",
	)

	if u.PasswordHash == "secret123" {
		t.Fatal("password must not be stored in plain text")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.PasswordHash),
		[]byte("secret123"),
	); err != nil {
		t.Errorf("hash does not match password: %v", err)
	}
}
