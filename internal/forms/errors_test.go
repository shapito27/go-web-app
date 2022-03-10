package forms

import "testing"

func TestErrors_Get(t *testing.T) {
	f := Form{
		Errors: errors{},
	}

	f.Errors["Name"] = append(f.Errors["Name"], "Name can't be empty")

	if f.Errors.Get("Name") != "Name can't be empty" {
		t.Error("We failed to get error")
	}

	if f.Errors.Get("Email") != "" {
		t.Error("We got value for not existing error")
	}
}

func TestErrors_Add(t *testing.T) {
	f := Form{
		Errors: errors{},
	}

	f.Errors.Add("Name", "Name can't be empty")

	if len(f.Errors["Name"]) == 0 {
		t.Error("We failed to add new error")
	}
}
