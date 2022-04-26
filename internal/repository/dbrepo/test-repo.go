package dbrepo

import (
	"database/sql"
	"time"

	"github.com/shapito27/go-web-app/internal/config"
	"github.com/shapito27/go-web-app/internal/models"
	"github.com/shapito27/go-web-app/internal/repository"
)

type testDBRepo struct {
	DB  *sql.DB
	App *config.AppConfig
}

func NewTestingDBRepo(app *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: app,
	}
}

// InsertReservation inserts a reservation to the table reservations
func (pr *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction insterts RoomRestriction to the table room_restrictions
func (pr *testDBRepo) InsertRoomRestriction(roomRestriction models.RoomRestriction) (int, error) {
	return 1, nil
}

// SearchAvailabilityByDatesAndRoomID returns true if availability exists for room id and false if no availability exists
func (pr *testDBRepo) SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms for given period
func (pr *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID returns room model by id
func (pr *testDBRepo) GetRoomByID(roomID int) (models.Room, error) {
	var room models.Room
	return room, nil
}
