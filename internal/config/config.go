package config

import (
	"text/template"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache       bool
	TemplatesCache map[string]*template.Template
	IsProduction   bool
	Session        *scs.SessionManager
}
