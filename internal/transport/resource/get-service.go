package resource

import (
	"fmt"
	"net/http"
)

func (a *APIHandler) GetService(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get service")
}
