package forms

import (
	"fmt"
	"net/http"
	"net/url"
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
