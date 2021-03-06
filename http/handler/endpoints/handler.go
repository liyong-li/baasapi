package endpoints

import (
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/http/proxy"
	"github.com/baasapi/baasapi/api/http/security"

	"net/http"

	"github.com/gorilla/mux"
)

const (
	// ErrEndpointManagementDisabled is an error raised when trying to access the endpoints management endpoints
	// when the server has been started with the --external-endpoints flag
	ErrEndpointManagementDisabled = baasapi.Error("Endpoint management is disabled")
)

func hideFields(endpoint *baasapi.Endpoint) {
	endpoint.AzureCredentials = baasapi.AzureCredentials{}
}

// Handler is the HTTP handler used to handle endpoint operations.
type Handler struct {
	*mux.Router
	authorizeEndpointManagement bool
	requestBouncer              *security.RequestBouncer
	EndpointService             baasapi.EndpointService
	EndpointGroupService        baasapi.EndpointGroupService
	FileService                 baasapi.FileService
	ProxyManager                *proxy.Manager
	Snapshotter                 baasapi.Snapshotter
	JobService                  baasapi.JobService
}

// NewHandler creates a handler to manage endpoint operations.
func NewHandler(bouncer *security.RequestBouncer, authorizeEndpointManagement bool) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
		authorizeEndpointManagement: authorizeEndpointManagement,
		requestBouncer:              bouncer,
	}

	h.Handle("/endpoints",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointCreate))).Methods(http.MethodPost)
	h.Handle("/endpoints/snapshot",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointSnapshots))).Methods(http.MethodPost)
	h.Handle("/endpoints",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.endpointList))).Methods(http.MethodGet)
	h.Handle("/endpoints/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.endpointInspect))).Methods(http.MethodGet)
	h.Handle("/endpoints/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointUpdate))).Methods(http.MethodPut)
	h.Handle("/endpoints/{id}/access",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointUpdateAccess))).Methods(http.MethodPut)
	h.Handle("/endpoints/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointDelete))).Methods(http.MethodDelete)
	h.Handle("/endpoints/{id}/extensions",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.endpointExtensionAdd))).Methods(http.MethodPost)
	h.Handle("/endpoints/{id}/extensions/{extensionType}",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.endpointExtensionRemove))).Methods(http.MethodDelete)
	h.Handle("/endpoints/{id}/job",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointJob))).Methods(http.MethodPost)
	h.Handle("/endpoints/{id}/snapshot",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.endpointSnapshot))).Methods(http.MethodPost)
	return h
}
