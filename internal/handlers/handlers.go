package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shapito27/go-web-app/internal/config"
	"github.com/shapito27/go-web-app/internal/forms"
	"github.com/shapito27/go-web-app/internal/helpers"
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
	render.RenderTemplate(w, r, "home", &models.TemplateData{})
}

// About page handler
func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "about", &models.TemplateData{})
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
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(out))
}

// Reservation handles requests to reservation form
func (rep *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation

	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles reservation post request
func (rep *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
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
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})

		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	rep.AppConfig.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Contact handles requests to contact page
func (rep *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact", &models.TemplateData{})
}

// ReservationSummary shows information from form
func (rep *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	reservation, ok := rep.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		errorMessage := "Can not pull reservation from the session"

		rep.AppConfig.ErrorLog.Println(errorMessage)
		rep.AppConfig.Session.Put(r.Context(), "error", errorMessage)

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	rep.AppConfig.Session.Remove(r.Context(), "reservation")

	data["reservation"] = reservation

	render.RenderTemplate(w, r, "reservation-summary", &models.TemplateData{
		Data: data,
	})
}
