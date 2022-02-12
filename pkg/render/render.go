package render

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/shapito27/go-web-app/pkg/config"
)

var functions = template.FuncMap{}

var appConfig *config.AppConfig

// Set config
func SetAppConfig(ac *config.AppConfig) {
	appConfig = ac
}

// To render templates
func RenderTemplate(w http.ResponseWriter, tmpl string) {
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

	err := t.Execute(w, nil)

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
