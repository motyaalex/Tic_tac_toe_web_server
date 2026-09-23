package main

import (
	"log"
	"net/http"
	"tic-tac-toe/internal/di"
	"tic-tac-toe/internal/web"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		di.Module,
		fx.Invoke(func(handler *web.Handler) {
			http.HandleFunc("/game/new", handler.CreateGame)
			http.Handle("/game/", handler)
			log.Println("Server starting on :8080")
			log.Fatal(http.ListenAndServe(":8080", nil))
		}),
	)

	app.Run()
}
