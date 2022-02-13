package handlers

import (
	"net/http"

	"github.com/shapito27/go-web-app/pkg/config"
	"github.com/shapito27/go-web-app/pkg/models"
	"github.com/shapito27/go-web-app/pkg/render"
)

var Repo *Repository

type Repository struct {
	AppConfig *config.AppConfig
}

// Create new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		AppConfig: a,
	}
}

// Set repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home page handler
func (rep *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "home", &models.TemplateData{})
}

// About page handler
func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Helo there"

	render.RenderTemplate(w, "about", &models.TemplateData{StringMap: stringMap})
}
