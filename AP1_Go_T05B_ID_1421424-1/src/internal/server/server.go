package server

import (
	"net/http"

	"tic-tac-toe/internal/config"
	"tic-tac-toe/internal/middleware"
	"tic-tac-toe/internal/web/handler"
)

func New(
	cfg config.Config,
	h *handler.GameHandler,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	authenticator *middleware.UserAuthenticator,
) *http.Server {
	// публичные маршруты
	publicMux := http.NewServeMux()
	publicMux.HandleFunc("POST /signup", authHandler.SignUp)
	publicMux.HandleFunc("POST /signin", authHandler.SignIn)
	publicMux.HandleFunc("POST /refresh-access-token", authHandler.RefreshAccessToken)
	publicMux.HandleFunc("POST /refresh-refresh-token", authHandler.RefreshRefreshToken)

	// защищённые маршруты
	protectedMux := http.NewServeMux()
	h.RegisterRoutes(protectedMux)
	userHandler.RegisterRoutes(protectedMux)
	protectedHandler := authenticator.Authenticate(protectedMux)

	// корневой mux
	rootMux := http.NewServeMux()
	rootMux.Handle("/signup", publicMux)
	rootMux.Handle("/signin", publicMux)
	rootMux.Handle("/", protectedHandler)
	rootMux.Handle("/refresh-access-token", publicMux)
	rootMux.Handle("/refresh-refresh-token", publicMux)

	var root http.Handler = rootMux
	root = middleware.Logging(root)
	root = middleware.Recovery(root)

	return &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      root,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
