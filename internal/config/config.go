package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/shapito27/go-web-app/internal/models"
)

type AppConfig struct {
	UseCache       bool
	TemplatesCache map[string]*template.Template
	IsProduction   bool
	Session        *scs.SessionManager
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	MailChan       chan models.MailData
}
