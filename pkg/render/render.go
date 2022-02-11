package render

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
)

var functions = template.FuncMap{}

// To render templates
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	templates, err := getTemplatesCache()
	if err != nil {
		fmt.Println("Error getting templates cache", err)
	}

	t, ok := templates[tmpl + ".page.tmpl"]
	if !ok {
		fmt.Println("Error getting template from cache")
	}
	
		buf := new(bytes.Buffer)

		_ = t.Execute(buf, nil)

		_, err = buf.WriteTo(w)

		if err != nil {
			fmt.Println("Error writing template to cachbrowser", err)
		}
	
}

func getTemplatesCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	//fmt.Println(" pages is", pages)

	matches, err := filepath.Glob("./templates/*.layout.tmpl")
	if err != nil {
		return myCache, err
	}
	//fmt.Println(" matches is", matches)

	for _, page := range pages {
		name := filepath.Base(page)

		//fmt.Println("Current page is", name)

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
