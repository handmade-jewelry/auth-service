package resource

import (
	"encoding/json"
	"net/http"

	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
)

func (a *APIHandler) GetServiceNameName(rw http.ResponseWriter, r *http.Request, name string) {
	srv, httpErr := a.serviceService.ServiceByName(r.Context(), name)
	if httpErr != nil {
		errors.WriteHTTPError(rw, httpErr)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(srv); err != nil {
		logger.Error("failed to encode response", err)
		errors.WriteHTTPError(rw, errors.InternalError())
		return
	}
}
