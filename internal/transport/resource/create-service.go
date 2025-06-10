package resource

import (
	"encoding/json"
	"net/http"

	"github.com/handmade-jewelry/auth-service/internal/service/service"
	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
)

func (a *APIHandler) PostService(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto service.ServiceDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		errors.WriteHTTPError(rw, errors.BadRequestError("Invalid request body"))

		return
	}

	if dto.Name == "" || dto.Host == "" {
		errors.WriteHTTPError(rw, errors.BadRequestError("Missing required fields"))
		return
	}

	srv, httpErr := a.serviceService.CreateService(r.Context(), &dto)
	if httpErr != nil {
		errors.WriteHTTPError(rw, httpErr)
		return
	}

	response := struct {
		ID int64 `json:"id"`
	}{
		ID: srv.ID,
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(rw).Encode(response); err != nil {
		logger.Error("failed to encode response", err)
		errors.WriteHTTPError(rw, errors.InternalError())
		return
	}
}
