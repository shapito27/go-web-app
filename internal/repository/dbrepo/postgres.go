package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/shapito27/go-web-app/internal/config"
	"github.com/shapito27/go-web-app/internal/models"
	"github.com/shapito27/go-web-app/internal/repository"
)

type postgresDBRepo struct {
	DB  *sql.DB
	App *config.AppConfig
}

func NewPostgresDBRepo(db *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		DB:  db,
		App: app,
	}
}

// InsertReservation inserts a reservation to the table reservations
func (pr *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, 
		room_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id;`

	err := pr.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction insterts RoomRestriction to the table room_restrictions
func (pr *postgresDBRepo) InsertRoomRestriction(roomRestriction models.RoomRestriction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id,  
		created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7) returning id;`

	err := pr.DB.QueryRowContext(ctx, stmt,
		roomRestriction.StartDate,
		roomRestriction.EndDate,
		roomRestriction.RoomID,
		roomRestriction.ReservationID,
		roomRestriction.RestrictionID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}
