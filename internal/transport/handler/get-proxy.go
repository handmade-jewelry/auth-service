package handler

import (
	"fmt"
	"net/http"
)

// todo rename?
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello from proxy")
}
