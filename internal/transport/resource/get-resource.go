package resource

import (
	"encoding/json"
	"net/http"
)

func (a *APIHandler) GetResourceId(w http.ResponseWriter, r *http.Request, id int) {
	resource, err := a.resourceService.Resource(r.Context(), int64(id))
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resource)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
