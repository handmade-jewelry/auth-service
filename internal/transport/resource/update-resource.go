package resource

import (
	"encoding/json"
	"net/http"

	"github.com/handmade-jewelry/auth-service/internal/service/resource"
	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
)

func (a *APIHandler) PutResourceId(rw http.ResponseWriter, r *http.Request, id int) {
	defer r.Body.Close()

	var dto resource.ResourceDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		errors.WriteHTTPError(rw, errors.BadRequestError("Invalid request body"))
		return
	}

	res, httpErr := a.resourceService.UpdateResource(r.Context(), dto, int64(id))
	if err != nil {
		errors.WriteHTTPError(rw, httpErr)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	resp := struct {
		ID int64 `json:"id"`
	}{
		ID: res.ID,
	}

	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		logger.Error("failed to encode response", err)
		errors.WriteHTTPError(rw, errors.InternalError())
		return
	}
}
