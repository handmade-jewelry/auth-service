package resource

import (
	"encoding/json"
	"net/http"
)

func (a *APIHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := a.userService.RoleList(r.Context())
	if err != nil {
		http.Error(w, "Failed to get roles", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(roles); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
