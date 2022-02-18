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
	remoteIp := r.RemoteAddr
	rep.AppConfig.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, "home", &models.TemplateData{})
}

// About page handler
func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Helo there"

	remoteIp := rep.AppConfig.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.RenderTemplate(w, "about", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (rep *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals", &models.TemplateData{})
}


func (rep *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors", &models.TemplateData{})
}


func (rep *Repository) Availablility(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availablility", &models.TemplateData{})
}

func (rep *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation", &models.TemplateData{})
}

func (rep *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact", &models.TemplateData{})
}