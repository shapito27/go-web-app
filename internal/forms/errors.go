package forms

type errors map[string][]string

// Add adds error message to given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get return first error message
func (e errors) Get(field string) string {
	if len(e[field]) == 0 {
		return ""
	}

	return e[field][0]
}
