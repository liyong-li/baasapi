package extensions

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/security"
)

// Handler is the HTTP handler used to handle extension operations.
type Handler struct {
	*mux.Router
	ExtensionService baasapi.ExtensionService
	ExtensionManager baasapi.ExtensionManager
}

// NewHandler creates a handler to manage extension operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}

	h.Handle("/extensions",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.extensionList))).Methods(http.MethodGet)
	h.Handle("/extensions",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.extensionCreate))).Methods(http.MethodPost)
	h.Handle("/extensions/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.extensionInspect))).Methods(http.MethodGet)
	h.Handle("/extensions/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.extensionDelete))).Methods(http.MethodDelete)
	h.Handle("/extensions/{id}/update",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.extensionUpdate))).Methods(http.MethodPost)

	return h
}
