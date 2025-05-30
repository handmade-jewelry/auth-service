package resource

import "github.com/handmade-jewelry/auth-service/internal/service/service"

type APIHandler struct {
	serviceService *service.Service
}

func NewAPIHandler(serviceService *service.Service) *APIHandler {
	return &APIHandler{
		serviceService: serviceService,
	}
}
