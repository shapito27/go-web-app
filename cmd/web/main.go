package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shapito27/go-web-app/internal/config"
	"github.com/shapito27/go-web-app/internal/driver"
	"github.com/shapito27/go-web-app/internal/handlers"
	"github.com/shapito27/go-web-app/internal/helpers"
	"github.com/shapito27/go-web-app/internal/models"
	"github.com/shapito27/go-web-app/internal/render"
)

//app listen this port
const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Printf("Starting application on http://localhost%s", portNumber)

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	//what I'm going to store in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// Setup environment
	app.IsProduction = false

	// Setup Loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	helpers.NewHelpers(&app)

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
		return nil, err
	}

	// save templates to config
	app.TemplatesCache = templates
	// false - development mode, true - production
	app.UseCache = false

	// connect to database
	log.Println("Connecting database...")
	dsn := "host=localhost port=5432 dbname=bookings user=postgres password=postgres"
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		fmt.Println("Error when connect sql", err)
		return nil, err
	}
	log.Println("Connected to database!")

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	// pass config to render package
	render.NewRenderer(&app)

	return db, nil
}
