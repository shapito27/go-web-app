package dbrepo

import (
	"database/sql"
	"errors"
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
	if res.RoomID == 2 {
		return 2, errors.New("Error")
	}
	return 1, nil
}

// InsertRoomRestriction insterts RoomRestriction to the table room_restrictions
func (pr *testDBRepo) InsertRoomRestriction(roomRestriction models.RoomRestriction) (int, error) {
	if roomRestriction.RoomID == 3 {
		return 2, errors.New("Error")
	}
	return 1, nil
}

// SearchAvailabilityByDatesAndRoomID returns true if availability exists for room id and false if no availability exists
func (pr *testDBRepo) SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms for given period
func (pr *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	layout := "2006-01-02"

	// error case when pass start 2022-06-26 and end 2022-06-29
	if start.Format(layout) == "2022-06-26" && end.Format(layout) == "2022-06-29" {
		return rooms, errors.New("Error")
	}

	// case when pass start 2022-06-27 and end 2022-06-30 - only one room available
	if start.Format(layout) == "2022-06-27" && end.Format(layout) == "2022-06-30" {
		rooms = append(rooms, models.Room{
			ID:       1,
			RoomName: "Nice room",
		})
	}

	return rooms, nil
}

// GetRoomByID returns room model by id
func (pr *testDBRepo) GetRoomByID(roomID int) (models.Room, error) {
	var room models.Room
	if roomID > 2 {
		return room, errors.New("Room not found")
	}
	return room, nil
}
