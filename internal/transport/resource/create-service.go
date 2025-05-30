package resource

import (
	"encoding/json"
	"github.com/handmade-jewelry/auth-service/internal/service/service"
	"net/http"
)

func (a *APIHandler) PostService(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto service.ServiceDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if dto.Name == "" || dto.Host == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	srv, err := a.serviceService.CreateService(r.Context(), &dto)
	if err != nil {
		http.Error(w, "Failed to create service", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID int64 `json:"id"`
	}{
		ID: srv.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
