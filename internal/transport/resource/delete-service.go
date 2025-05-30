package resource

import (
	"net/http"
)

func (a *APIHandler) DeleteServiceId(w http.ResponseWriter, r *http.Request, id int) {
	err := a.serviceService.DeleteService(r.Context(), int64(id))
	if err != nil {
		http.Error(w, "Failed to get service", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
