package resource

import (
	"encoding/json"
	"net/http"
)

func (a *APIHandler) GetServiceNameName(w http.ResponseWriter, r *http.Request, name string) {
	srv, err := a.serviceService.ServiceByName(r.Context(), name)
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
