package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shapito27/go-web-app/internal/config"
	"github.com/shapito27/go-web-app/internal/driver"
	"github.com/shapito27/go-web-app/internal/forms"
	"github.com/shapito27/go-web-app/internal/helpers"
	"github.com/shapito27/go-web-app/internal/models"
	"github.com/shapito27/go-web-app/internal/render"
	"github.com/shapito27/go-web-app/internal/repository"
	"github.com/shapito27/go-web-app/internal/repository/dbrepo"
)

var Repo *Repository

type Repository struct {
	AppConfig *config.AppConfig
	DB        repository.DatabaseRepo
}

// Create new Repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		AppConfig: a,
		DB:        dbrepo.NewPostgresDBRepo(db.SQL, a),
	}
}

// Create new Repository
func NewTestingRepo(a *config.AppConfig) *Repository {
	return &Repository{
		AppConfig: a,
		DB:        dbrepo.NewTestingDBRepo(a),
	}
}

// Set repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home page handler
func (rep *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home", &models.TemplateData{})
}

// About page handler
func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about", &models.TemplateData{})
}

func (rep *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals", &models.TemplateData{})
}

func (rep *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors", &models.TemplateData{})
}

func (rep *Repository) Availablility(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availablility", &models.TemplateData{})
}

func (rep *Repository) PostAvailablility(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {

		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := rep.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		rep.AppConfig.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	for _, room := range rooms {
		rep.AppConfig.InfoLog.Println(room)
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	rep.AppConfig.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room", &models.TemplateData{
		Data: data,
	})
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
	res, ok := rep.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := rep.DB.GetRoomByID(res.RoomID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room = room

	data := make(map[string]interface{})
	data["reservation"] = res

	rep.AppConfig.Session.Put(r.Context(), "reservation", res)

	stringMap := make(map[string]string)
	stringMap["start_date"] = res.StartDate.Format("2006-01-02")
	stringMap["end_date"] = res.EndDate.Format("2006-01-02")

	render.Template(w, r, "make-reservation", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles reservation post request
func (rep *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation, ok := rep.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})

		data["reservation"] = reservation
		http.Error(w, "Validation failed", http.StatusSeeOther)
		render.Template(w, r, "make-reservation", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	// save to DB
	newReservationID, err := rep.DB.InsertReservation(reservation)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 2,
	}

	_, err = rep.DB.InsertRoomRestriction(restriction)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// send notifications to guest
	guestMessage := `Hi, %s<br>
This is confirtm your reservation from %s to %s.
`
	mail := models.MailData{
		From:    "me@gmail.com",
		To:      reservation.Email,
		Subject: "Reservation confirmation",
		Content: fmt.Sprintf(guestMessage, reservation.FirstName, reservation.StartDate, reservation.EndDate),
	}
	rep.AppConfig.MailChan <- mail

	// send notifications to owner
	guestMessage = `New reservation<br>
%s, made a reservation from %s to %s.
`
	mail = models.MailData{
		From:    "me@gmail.com",
		To:      "me@gmail.com",
		Subject: "Reservation confirmation",
		Content: fmt.Sprintf(guestMessage, reservation.FirstName, reservation.StartDate, reservation.EndDate),
	}
	rep.AppConfig.MailChan <- mail

	rep.AppConfig.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Contact handles requests to contact page
func (rep *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact", &models.TemplateData{})
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

	stringMap := make(map[string]string)
	stringMap["start_date"] = reservation.StartDate.Format("2006-01-02")
	stringMap["end_date"] = reservation.EndDate.Format("2006-01-02")

	render.Template(w, r, "reservation-summary", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom save room id to session and redirect to make reservation date
func (rep *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := rep.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}
	res.RoomID = roomId
	rep.AppConfig.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
