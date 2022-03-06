package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var mh myHandler

	noSurfHandler := NoSurf(&mh)

	switch noSurfHandler.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("Wrong type. NoSurf has to return http.Handler, but it's %T", noSurfHandler))
	}
}

func TestLoadSession(t *testing.T) {
	var mh myHandler

	ls := LoadSession(&mh)

	switch ls.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("Wrong type. LoadSession has to return http.Handler, but it's %T", ls))
	}
}
