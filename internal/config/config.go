package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache       bool
	TemplatesCache map[string]*template.Template
	IsProduction   bool
	Session        *scs.SessionManager
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
}
