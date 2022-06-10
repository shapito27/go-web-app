package main

import (
	"encoding/gob"
	"flag"
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

	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

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

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "postgres", "Database pass")
	dbSSL := flag.String("dbssl", "disable", "Database ssl")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		println("Missing required flags")
		os.Exit(1)
	}

	// create chanel for emails
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// Setup environment
	app.IsProduction = *inProduction

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
	app.UseCache = *useCache

	// connect to database
	log.Println("Connecting database...")
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
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
