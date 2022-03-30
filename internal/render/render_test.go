package render

import (
	"net/http"
	"testing"

	"github.com/shapito27/go-web-app/internal/config"
	"github.com/shapito27/go-web-app/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	templateData := &models.TemplateData{}

	r, err := buildTestRequest()
	if err != nil {
		t.Error("can't build request", err)
	}

	flashMessage := "test_flash"
	session.Put(r.Context(), "flash", flashMessage)

	templateData = addDefaultData(templateData, r)

	if templateData.Flash != flashMessage {
		t.Error("AddDefaultData failed to set flash message")
	}
}

func buildTestRequest() (*http.Request, error) {
	// building request
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}

	// prepare request context using session data
	ctx, _ := session.Load(r.Context(), r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplate = "./../../templates"

	templates, err := GetTemplatesCache()

	if err != nil {
		t.Error(err)

	}

	appConfig.TemplatesCache = templates

	w := &myResponseWriter{}
	r, err := buildTestRequest()
	if err != nil {
		t.Error(err)

	}

	templateData := &models.TemplateData{}

	err = Template(w, r, "home", templateData)
	if err != nil {
		t.Error("Can not render existed template", err)
	}

	err = Template(w, r, "6666", templateData)
	if err == nil {
		t.Error("Got not existed template")
	}
}

func TestSetAppConfig(t *testing.T) {

	NewRenderer(&config.AppConfig{
		UseCache:     true,
		IsProduction: true,
	})
	if appConfig.UseCache != true {
		t.Error("SetAppConfig doesn't set UseCache correctly ")
	}

	if appConfig.IsProduction != true {
		t.Error("SetAppConfig doesn't set IsProduction correctly ")
	}
}

func TestGetTemplatesCache(t *testing.T) {
	pathToTemplate = "./../../templates"
	_, err := GetTemplatesCache()
	if err != nil {
		t.Error(err)
	}
}
