package forms

import (
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	err := r.ParseForm()
	if err != nil {
		t.Error(err)
	}

	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("Form without error, but got 'no valid'")
	}
}

func TestNewEmptyForm(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	err := r.ParseForm()
	if err != nil {
		t.Error(err)
	}

	f := New(r.PostForm)

	for k, v := range r.PostForm {
		if !reflect.DeepEqual(f.Values[k], v) {
			t.Error("Values from request are not set")
		}
	}
}

func TestNewNotEmptyForm(t *testing.T) {
	values := url.Values{}
	values.Add("Name", "David")
	values.Add("Email", "david@gmail.com")

	bodyReader := strings.NewReader(values.Encode())

	r := httptest.NewRequest("POST", "/", bodyReader)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := r.ParseForm()
	if err != nil {
		t.Error(err)
	}

	f := New(r.PostForm)

	for k, v := range r.PostForm {
		if !reflect.DeepEqual(f.Values[k], v) {
			t.Error("Values from request are not set")
		}
	}
}

func TestForm_Has(t *testing.T) {
	values := url.Values{}
	values.Add("Name", "David")

	bodyReader := strings.NewReader(values.Encode())

	r := httptest.NewRequest("POST", "/", bodyReader)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := r.ParseForm()
	if err != nil {
		t.Error(err)
	}

	f := New(r.PostForm)

	if !f.Has("Name") {
		t.Error("Form has to have Name field")
	}

	if f.Has("Phone") {
		t.Error("Form hasn't to have Phone field")
	}
}

func TestForm_Required(t *testing.T) {
	values := url.Values{}
	values.Add("Name", "David")
	values.Add("Email", "")

	bodyReader := strings.NewReader(values.Encode())

	r := httptest.NewRequest("POST", "/", bodyReader)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := r.ParseForm()
	if err != nil {
		t.Error(err)
	}

	f := New(r.PostForm)

	f.Required("Name")
	f.Required("Email")

	if len(f.Errors["Name"]) != 0 {
		t.Error("Name field is required and it set but got new entry in Errors")
	}

	if len(f.Errors["Email"]) == 0 {
		t.Error("There is required field is empty, but no error in the Form")
	}
}

func TestForm_MinLength(t *testing.T) {
	values := url.Values{}
	values.Add("Name", "David")
	values.Add("LastName", "Y")

	bodyReader := strings.NewReader(values.Encode())

	r := httptest.NewRequest("POST", "/", bodyReader)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := r.ParseForm()
	if err != nil {
		t.Error(err)
	}

	f := New(r.PostForm)

	f.MinLength("Name", 3)
	f.MinLength("LastName", 3)

	if len(f.Errors["Name"]) != 0 {
		t.Error("Name field is long enough but got new entry in Form.Errors")
	}

	if len(f.Errors["LastName"]) == 0 {
		t.Error("There is not long enough LastName, but no error in the Form.Errors")
	}
}

// TestForm_IsEmail1 checks that there is hasn't to be errors if input valid email
func TestForm_IsEmail1(t *testing.T) {
	values := url.Values{}
	values.Add("Email", "david@gmail.com")

	bodyReader := strings.NewReader(values.Encode())

	r := httptest.NewRequest("POST", "/", bodyReader)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := r.ParseForm()
	if err != nil {
		t.Error(err)
	}

	f := New(r.PostForm)

	f.IsEmail("Email")

	if len(f.Errors["Email"]) != 0 {
		t.Error("Email field is valid but got new entry in Form.Errors")
	}
}

// TestForm_IsEmail2 checks that there is has to be errors if input not valid email
func TestForm_IsEmail2(t *testing.T) {
	values := url.Values{}
	values.Add("Email", "8-800-555-35-35")

	bodyReader := strings.NewReader(values.Encode())

	r := httptest.NewRequest("POST", "/", bodyReader)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := r.ParseForm()
	if err != nil {
		t.Error(err)
	}

	f := New(r.PostForm)

	f.IsEmail("Email")

	if len(f.Errors["Email"]) == 0 {
		t.Error("Email field is not valid but got new entry in Form.Errors")
	}
}
