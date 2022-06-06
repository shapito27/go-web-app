package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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

// TestRepository_PostReservation
func TestRepository_PostReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Red room",
		},
	}

	strBody := "first_name=Rusaln&last_name=Jora&email=ll@ll.ru&phone=88005553535&start_date=2023-01-01&end_date=2024-01-01"
	body := strings.NewReader(strBody)

	req, _ := http.NewRequest("POST", "/make-reservation", body)
	ctx := getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Handler PostReservation returned wrong response code: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// case when wrong body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostReservation returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case when no session
	strBody = "first_name=Rusaln&last_name=Jora&email=ll@ll.ru&phone=88005553535&start_date=2023-01-01&end_date=2024-01-01"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostReservation returned wrong response code when no session: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case when not valid data
	strBody = "first_name=R&last_name=Jora&email=ll@ll.ru&phone=88005553535&start_date=2023-01-01&end_date=2024-01-01"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Handler PostReservation returned wrong response code when invalid name: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// case when insert reservation fails
	strBody = "first_name=Ruslan&last_name=Jora&email=ll@ll.ru&phone=88005553535&start_date=2023-01-01&end_date=2024-01-01"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	reservation.RoomID = 2
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostReservation returned wrong response code when insert reservation fails: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case when insert RoomRestriction fails
	strBody = "first_name=Ruslan&last_name=Jora&email=ll@ll.ru&phone=88005553535&start_date=2023-01-01&end_date=2024-01-01"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/make-reservation", body)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	reservation.RoomID = 3
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostReservation returned wrong response code when insert RoomRestriction fails: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

// TestRepository_PostAvailablility
func TestRepository_PostAvailablility(t *testing.T) {
	strBody := "start=2022-06-27&end=2022-06-30"
	body := strings.NewReader(strBody)

	req, _ := http.NewRequest("POST", "/search-availability", body)
	ctx := getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailablility)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler PostAvailablility returned wrong response code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// cace trying fail SearchAvailabilityForAllRooms
	strBody = "start=2022-06-26&end=2022-06-29"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/search-availability", body)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailablility)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostAvailablility returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case when no post body
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailablility)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostAvailablility returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case when no start date
	strBody = "end=2022-06-22"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/search-availability", body)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailablility)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostAvailablility returned wrong response code when start empty: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case when no end date
	strBody = "start=2022-06-22"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/search-availability", body)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailablility)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostAvailablility returned wrong response code when end empty: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case when no rooms for dates
	strBody = "start=2022-06-21&end=2022-06-22"
	body = strings.NewReader(strBody)

	req, _ = http.NewRequest("POST", "/search-availability", body)
	ctx = getContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailablility)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Handler PostAvailablility returned wrong response code when no rooms for dates: got %d, expected %d", rr.Code, http.StatusSeeOther)
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
