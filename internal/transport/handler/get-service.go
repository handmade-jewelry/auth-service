package handler

import (
	"fmt"
	"net/http"
)

func (a *APIHandler) GetAdminService(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get service")
}
