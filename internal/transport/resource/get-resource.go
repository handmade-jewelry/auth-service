package resource

import (
	"encoding/json"
	"net/http"

	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
)

func (a *APIHandler) GetResourceId(rw http.ResponseWriter, req *http.Request, id int) {
	resource, httpErr := a.resourceService.Resource(req.Context(), int64(id))
	if httpErr != nil {
		errors.WriteHTTPError(rw, httpErr)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(resource)
	if err != nil {
		logger.Error("failed to encode response", err)
		errors.WriteHTTPError(rw, errors.InternalError())
		return
	}
}
