package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shapito27/go-web-app/internal/config"
	"github.com/shapito27/go-web-app/internal/helpers"
	"github.com/shapito27/go-web-app/internal/models"
)

var session *scs.SessionManager
var app config.AppConfig

func TestMain(m *testing.M) {
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

	appConfig = &app

	os.Exit(m.Run())
}

type myResponseWriter struct{}

func (rw *myResponseWriter) Header() http.Header {
	return http.Header{}
}

func (rw *myResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (rw *myResponseWriter) WriteHeader(statusCode int) {

}
