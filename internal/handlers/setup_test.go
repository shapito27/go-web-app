package handlers

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/shapito27/go-web-app/internal/config"
	"github.com/shapito27/go-web-app/internal/models"
	"github.com/shapito27/go-web-app/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager

var functions = template.FuncMap{}

var pathToTemplate = "./../../templates"

func getRoutes() http.Handler {
	//what I'm going to store in the session
	gob.Register(models.Reservation{})

	// Setup environment
	app.IsProduction = false

	// Setup session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	// getting templates list
	templates, err := GetTestTemplatesCache()
	if err != nil {
		fmt.Println("Error getting templates cache", err)
	}

	// save templates to config
	app.TemplatesCache = templates
	// false - development mode, true - production
	app.UseCache = true

	repo := NewRepo(&app)
	NewHandlers(repo)

	// pass config to render package
	render.SetAppConfig(&app)

	mux := chi.NewRouter()

	// Middlewares
	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	mux.Use(LoadSession)

	// Routes
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availablility)
	mux.Post("/search-availability", Repo.PostAvailablility)
	mux.Post("/search-availability-json", Repo.PostAvailablilityJSON)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf sets CSRF token for every request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.IsProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// LoadSession loads and saves session on every request
func LoadSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

//collect all templates then merge them with layout
func GetTestTemplatesCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplate))
	if err != nil {
		return myCache, err
	}

	matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
