package main

import (
	"context"
	"net/http"

	"go.uber.org/fx"

	"tic-tac-toe/internal/di"
)

func main() {
	fx.New(
		di.Module,
		fx.Invoke(func(lc fx.Lifecycle, s *http.Server) {
			lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
				go func() { _ = s.ListenAndServe() }()
				return nil
			},
				OnStop: func(ctx context.Context) error {
					return s.Shutdown(ctx)
				}})
		}),
	).Run()
}
