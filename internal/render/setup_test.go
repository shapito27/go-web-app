package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shapito27/go-web-app/internal/config"
	"github.com/shapito27/go-web-app/internal/models"
)

var session *scs.SessionManager
var app config.AppConfig

func TestMain(m *testing.M) {
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

	appConfig = &app

	os.Exit(m.Run())
}
