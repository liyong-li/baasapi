package resourcecontrols

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// Handler is the HTTP handler used to handle resource control operations.
type Handler struct {
	*mux.Router
	ResourceControlService baasapi.ResourceControlService
}

// NewHandler creates a handler to manage resource control operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/resource_controls",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.resourceControlCreate))).Methods(http.MethodPost)
	h.Handle("/resource_controls/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.resourceControlUpdate))).Methods(http.MethodPut)
	h.Handle("/resource_controls/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.resourceControlDelete))).Methods(http.MethodDelete)

	return h
}
