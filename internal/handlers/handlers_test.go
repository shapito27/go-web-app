package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shapito27/go-web-app/internal/models"
)

type PostData struct {
	key   string
	value string
}

var tests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
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

// TestRepository_Reservation test cases related to making reservation
func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Red room",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getContext(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler Reservation returned wrong response code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// test case where Reservation is not in the session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler Reservation returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case with non existing room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	reservation.RoomID = 3
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler Reservation returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

// getContext return context by request
func getContext(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
