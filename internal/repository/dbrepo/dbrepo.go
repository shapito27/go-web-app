package dbrepo

import (
	"database/sql"

	"github.com/shapito27/go-web-app/internal/config"
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
