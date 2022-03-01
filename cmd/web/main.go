package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shapito27/go-web-app/internal/config"
	"github.com/shapito27/go-web-app/internal/handlers"
	"github.com/shapito27/go-web-app/internal/models"
	"github.com/shapito27/go-web-app/internal/render"
)

//app listen this port
const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	//what I'm going to store in the session
	gob.Register(models.Reservation{})

	// Setup environment
	app.IsProduction = false

	// Setup session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	// getting templates list
	templates, err := render.GetTemplatesCache()
	if err != nil {
		fmt.Println("Error getting templates cache", err)
	}

	// save templates to config
	app.TemplatesCache = templates
	// false - development mode, true - production
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// pass config to render package
	render.SetAppConfig(&app)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	err = srv.ListenAndServe()
	log.Fatal(err)
}
