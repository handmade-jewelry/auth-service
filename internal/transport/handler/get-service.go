package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) GetAdminService(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get service")
}
