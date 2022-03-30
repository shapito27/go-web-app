package repository

import "github.com/shapito27/go-web-app/internal/models"

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) error
}
