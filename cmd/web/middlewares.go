package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// just test middleware
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf("Hit the page: %s", r.URL))

		next.ServeHTTP(w, r)
	})
}

// NoSurf sets CSRF token for every request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.IsProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// LoadSession loads and saves session on every request
func LoadSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
