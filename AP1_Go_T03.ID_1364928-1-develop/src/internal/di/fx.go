package di

import (
	"tic-tac-toe/internal/application"
	"tic-tac-toe/internal/datasource"
	"tic-tac-toe/internal/domain"
	"tic-tac-toe/internal/web"

	"go.uber.org/fx"
)

var Module = fx.Options( // Проводка/слейка между всеми слоями
	fx.Provide(
		datasource.NewGameStorage,    // singleton  Инициализация и создание хранилища
		datasource.NewGameRepository, // Определяем репозиторий для хранилища *GameStorage
		application.NewGameService,   // сервис для domain.GameRepository и правила игры, алгоритм Минимакс
		web.NewHandler,               // хендлер веб слоя, принимает и отдает запросы (фронтенд), вызывает методы игры
	),
	fx.Invoke(func(domain.GameRepository) {}), // гарантируем инициализацию сервера
)
