package resource

import "net/http"

func (a *APIHandler) DeleteResourceId(w http.ResponseWriter, r *http.Request, id int) {
	err := a.resourceService.DeleteResource(r.Context(), int64(id))
	if err != nil {
		http.Error(w, "Failed to delete resource", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
