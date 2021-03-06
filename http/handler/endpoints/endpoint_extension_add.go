package endpoints

// TODO: legacy extension management

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type endpointExtensionAddPayload struct {
	Type int
	URL  string
}

func (payload *endpointExtensionAddPayload) Validate(r *http.Request) error {
	if payload.Type != 1 {
		return baasapi.Error("Invalid type value. Value must be one of: 1 (Storidge)")
	}
	if payload.Type == 1 && govalidator.IsNull(payload.URL) {
		return baasapi.Error("Invalid extension URL")
	}
	return nil
}

// POST request on /api/endpoints/:id/extensions
func (handler *Handler) endpointExtensionAdd(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	endpointID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid endpoint identifier route variable", err}
	}

	endpoint, err := handler.EndpointService.Endpoint(baasapi.EndpointID(endpointID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find an endpoint with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an endpoint with the specified identifier inside the database", err}
	}

	var payload endpointExtensionAddPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	extensionType := baasapi.EndpointExtensionType(payload.Type)

	var extension *baasapi.EndpointExtension
	for _, ext := range endpoint.Extensions {
		if ext.Type == extensionType {
			extension = &ext
		}
	}

	if extension != nil {
		extension.URL = payload.URL
	} else {
		extension = &baasapi.EndpointExtension{
			Type: extensionType,
			URL:  payload.URL,
		}
		endpoint.Extensions = append(endpoint.Extensions, *extension)
	}

	err = handler.EndpointService.UpdateEndpoint(endpoint.ID, endpoint)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist endpoint changes inside the database", err}
	}

	return response.JSON(w, extension)
}
