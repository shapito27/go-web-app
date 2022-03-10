package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
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
func (f *Form) Has(field string) bool {
	getField := f.Values.Get(field)
	if getField == "" {
		f.Errors.Add(field, fmt.Sprintf("%s can't be blank", field))
		return false
	}

	return true
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
func (f *Form) MinLength(field string, length int) {
	value := f.Values.Get(field)
	if len(value) < length {
		f.Errors.Add(field, fmt.Sprintf("%s can't be less than %d symbols", field, length))
	}
}

// IsEmail checks email valid
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Values.Get(field)) {
		f.Errors.Add(field, fmt.Sprintf("%s is not valid", field))
	}
}
