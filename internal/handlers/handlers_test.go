package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type PostData struct {
	key   string
	value string
}

var tests = []struct {
	name               string
	url                string
	method             string
	params             []PostData
	expectedStatusCode int
}{
	{"home", "/", "GET", []PostData{}, http.StatusOK},
	{"about", "/about", "GET", []PostData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", []PostData{}, http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", []PostData{}, http.StatusOK},

	{"search-availability", "/search-availability", "GET", []PostData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", []PostData{}, http.StatusOK},
	{"contact", "/contact", "GET", []PostData{}, http.StatusOK},
	//{"reservation-summary", "/reservation-summary", "GET", []PostData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range tests {
		if test.method == "GET" {
			response, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Error(err)
			}
			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("Not correct status code. For %s expected %d but got %d", test.url, test.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
