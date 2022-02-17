package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shapito27/go-web-app/pkg/config"
	"github.com/shapito27/go-web-app/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// Middlewares
	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(LoadSession)

	// Routes
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
