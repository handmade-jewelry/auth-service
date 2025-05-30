package resource

import (
	"encoding/json"
	"net/http"
)

func (a *APIHandler) GetServiceId(w http.ResponseWriter, r *http.Request, id int) {
	srv, err := a.serviceService.ServiceByID(r.Context(), int64(id))
	if err != nil {
		http.Error(w, "Failed to get service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(srv); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
