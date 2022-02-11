package render

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
)

var functions = template.FuncMap{}

// To render templates
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	templates, err := getTemplatesList()
	if err != nil {
		fmt.Println("Error getting templates cache", err)
	}

	t, ok := templates[tmpl+".page.tmpl"]
	if !ok {
		fmt.Println("Error getting template from cache")
	}

	err = t.Execute(w, nil)

	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}

}

//collect all templates then merge them with layout
func getTemplatesList() (map[string]*template.Template, error) {
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
