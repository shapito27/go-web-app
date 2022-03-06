package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/shapito27/go-web-app/internal/config"
)

func TestRoutes(t *testing.T) {
	var appConfig config.AppConfig

	handler := routes(&appConfig)

	switch handler.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("Wrong type. routes has to return http.Handler, but it's %T", handler))
	}
}
