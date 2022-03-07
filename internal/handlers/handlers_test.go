package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	{"search-availability", "/search-availability", "POST", []PostData{
		{key: "start", value: "2022-03-01"},
		{key: "end", value: "2022-03-07"},
	}, http.StatusOK},
	{"search-availability-json", "/search-availability-json", "POST", []PostData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "POST", []PostData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Toe"},
		{key: "phone", value: "8-800-888"},
		{key: "email", value: "test@yan.ru"},
	}, http.StatusOK},
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
		} else if test.method == "POST" {
			values := url.Values{}

			for _, param := range test.params {
				values.Add(param.key, param.value)
			}

			response, err := testServer.Client().PostForm(testServer.URL + test.url, values)
			if err != nil {
				t.Error(err)
			}
			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("Not correct status code. For %s expected %d but got %d", test.url, test.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
