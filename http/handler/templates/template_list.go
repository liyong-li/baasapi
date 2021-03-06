package templates

import (
	"encoding/json"
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/client"
	"github.com/baasapi/baasapi/api/http/security"
)

// GET request on /api/templates
func (handler *Handler) templateList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	settings, err := handler.SettingsService.Settings()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve settings from the database", err}
	}

	var templates []baasapi.Template
	if settings.TemplatesURL == "" {
		templates, err = handler.TemplateService.Templates()
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve templates from the database", err}
		}
	} else {
		var templateData []byte
		templateData, err = client.Get(settings.TemplatesURL, 0)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve external templates", err}
		}

		err = json.Unmarshal(templateData, &templates)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to parse external templates", err}
		}
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	filteredTemplates := security.FilterTemplates(templates, securityContext)
	return response.JSON(w, filteredTemplates)
}
