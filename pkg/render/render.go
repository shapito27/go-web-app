package render

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	"github.com/shapito27/go-web-app/pkg/config"
	"github.com/shapito27/go-web-app/pkg/models"
)

var functions = template.FuncMap{}

var appConfig *config.AppConfig

// Set config
func SetAppConfig(ac *config.AppConfig) {
	appConfig = ac
}

// add default data to all template data
func addDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {
	templateData.CSRFToken = nosurf.Token(r)

	return templateData
}

// To render templates
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) {
	var templates map[string]*template.Template

	if appConfig.UseCache {
		// getting templates list from config
		templates = appConfig.TemplatesCache
	} else {
		var err error
		templates, err = GetTemplatesCache()
		if err != nil {
			fmt.Println("Error getting template from cache", err)
		}
	}

	t, ok := templates[tmpl+".page.tmpl"]
	if !ok {
		fmt.Println("Error getting template from cache")
	}

	data = addDefaultData(data, r)

	err := t.Execute(w, data)

	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}

}

//collect all templates then merge them with layout
func GetTemplatesCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	matches, err := filepath.Glob("./templates/*.layout.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
