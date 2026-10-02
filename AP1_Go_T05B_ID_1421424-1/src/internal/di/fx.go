package di

import (
	"tic-tac-toe/internal/application/service"
	"tic-tac-toe/internal/config"
	"tic-tac-toe/internal/datasource/postgres"
	"tic-tac-toe/internal/datasource/repository"
	"tic-tac-toe/internal/middleware"
	"tic-tac-toe/internal/server"
	"tic-tac-toe/internal/web/handler"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		config.Load,
		postgres.NewPool,
		repository.NewGameRepository,
		repository.NewUserRepository,
		service.NewGameService,
		service.NewUserService,
		service.NewJwtProvider,
		service.NewAuthService,
		handler.NewGameHandler,
		handler.NewUserHandler,
		handler.NewAuthHandler,
		middleware.NewUserAuthenticator,
		server.New,
	),
)
