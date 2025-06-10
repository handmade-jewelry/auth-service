package resource

import (
	"net/http"

	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
)

func (a *APIHandler) DeleteResourceId(rw http.ResponseWriter, req *http.Request, id int) {
	httpErr := a.resourceService.DeleteResource(req.Context(), int64(id))
	if httpErr != nil {
		errors.WriteHTTPError(rw, httpErr)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
