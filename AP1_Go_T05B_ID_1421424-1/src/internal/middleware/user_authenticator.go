package middleware

import (
	"context"
	"net/http"
	"strings"

	"tic-tac-toe/internal/application/service"

	"github.com/google/uuid"
)

type ctxKey string

const userIDKey ctxKey = "userID"

type UserAuthenticator struct {
	jwt *service.JwtProvider
}

func NewUserAuthenticator(jwt *service.JwtProvider) *UserAuthenticator { // Новая авторизация через jwt
	return &UserAuthenticator{
		jwt: jwt,
	}
}

// Authenticate — middleware, требующая валидный Basic Auth. // Переделано через jwt
func (a *UserAuthenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken := parts[1]

		userID, err := a.jwt.ValidateAccessToken(accessToken)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// прокидываем UUID пользователя в context для хендлеров
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}
