package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shapito27/go-web-app/pkg/config"
	"github.com/shapito27/go-web-app/pkg/handlers"
	"github.com/shapito27/go-web-app/pkg/render"
)

//app listen this port
const portNumber = ":8080"

func main() {
	var config config.AppConfig

	// getting templates list
	templates, err := render.GetTemplatesCache()
	if err != nil {
		fmt.Println("Error getting templates cache", err)
	}

	// save templates to config
	config.TemplatesCache = templates
	// false - development mode, true - production
	config.UseCache = true

	repo := handlers.NewRepo(&config)
	handlers.NewHandlers(repo)

	// pass config to render package
	render.SetAppConfig(&config)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&config),
	}

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	err = srv.ListenAndServe()
	log.Fatal(err)
}
