package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shapito27/go-web-app/internal/config"
	"github.com/shapito27/go-web-app/internal/forms"
	"github.com/shapito27/go-web-app/internal/models"
	"github.com/shapito27/go-web-app/internal/render"
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

	render.RenderTemplate(w, r, "home", &models.TemplateData{})
}

// About page handler
func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Helo there"

	remoteIp := rep.AppConfig.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.RenderTemplate(w, r, "about", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (rep *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals", &models.TemplateData{})
}

func (rep *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors", &models.TemplateData{})
}

func (rep *Repository) Availablility(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availablility", &models.TemplateData{})
}

func (rep *Repository) PostAvailablility(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("dates is %s - %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (rep *Repository) PostAvailablilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "    ")

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(out))
}

func (rep *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation

	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (rep *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)

	if !form.Valid() {
		data := make(map[string]interface{})

		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}
}

func (rep *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact", &models.TemplateData{})
}
