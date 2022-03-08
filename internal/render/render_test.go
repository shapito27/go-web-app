package render

import (
	"net/http"
	"testing"

	"github.com/shapito27/go-web-app/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	templateData := &models.TemplateData{}

	// building request
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error("Failed to create request")
	}

	// prepare request context using session data
	//ctx := r.Context()
	ctx, _ := session.Load(r.Context(), r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	flashMessage := "test_flash"
	session.Put(r.Context(), "flash", flashMessage)

	templateData = addDefaultData(templateData, r)

	if templateData.Flash != flashMessage {
		t.Error("AddDefaultData failed to set flash message")
	}
}
