package config

import "text/template"

type AppConfig struct {
	UseCache      bool
	TemplatesCache map[string]*template.Template
}
