package resource

import (
	"net/http"

	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
)

func (a *APIHandler) DeleteServiceId(rw http.ResponseWriter, req *http.Request, id int) {
	httpErr := a.serviceService.DeleteService(req.Context(), int64(id))
	if httpErr != nil {
		errors.WriteHTTPError(rw, httpErr)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
