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

// SearchAvailabilityByDatesAndRoomID returns true if availability exists for room id and false if no availability exists
func (pr *postgresDBRepo) SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error) {
	var numRows int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		count(id)
	from 
		room_restrictions
	where 
	room_id = $1
		$2 < end_date and $3 > start_date;`

	row := pr.DB.QueryRowContext(ctx, query, roomID, end, start)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms for given period
func (pr *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
	select 
		r.id, r.room_name
	from
		rooms r
	where 
		r.id not in (
			select
				room_id
			from
				room_restrictions rr
			where
				$1 < rr.end_date and $2 > rr.start_date 
	)
	`

	rows, err := pr.DB.QueryContext(ctx, query, end, start)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room models.Room

		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID returns room model by id
func (pr *postgresDBRepo) GetRoomByID(roomID int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var room models.Room

	query := `
	select id, room_name, created_at, updated_at 
	from rooms
	where id=$1
	`

	row := pr.DB.QueryRowContext(ctx, query, roomID)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}

	return room, nil
}
