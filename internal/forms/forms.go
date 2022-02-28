package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Form stores form values and errors
type Form struct {
	Values url.Values
	Errors errors
}

// Valid return true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initilizes Form data
func New(data url.Values) *Form {
	return &Form{
		data,
		//errors(map[string][]string{}),
		errors{},
	}
}

// Has checks if request has field
func (f *Form) Has(field string, r *http.Request) bool {
	getField := r.Form.Get(field)
	if getField == "" {
		f.Errors.Add(field, fmt.Sprintf("%s can't be blank", field))
		return false
	}

	return r.Form.Has(field)
}

// Required checks all given fields have value
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Values.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("%s can't be blank", field))
		}
	}
}

// MinLength checks field lenght not less then given length
func (f *Form) MinLength(field string, length int, r *http.Request) {
	value := r.Form.Get(field)
	if len(value) < length {
		f.Errors.Add(field, fmt.Sprintf("%s can't be less than %d symbols", field, length))
	}
}
