package resource

import (
	"encoding/json"
	"github.com/handmade-jewelry/auth-service/internal/service/resource"
	"net/http"
)

func (a *APIHandler) PostResource(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto resource.ResourceDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	res, err := a.resourceService.CreateResource(r.Context(), dto)
	if err != nil {
		http.Error(w, "Failed to create resource: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := struct {
		ID int64 `json:"id"`
	}{
		ID: res.ID,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
