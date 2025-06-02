package resource

import (
	"github.com/handmade-jewelry/auth-service/internal/service/resource"
	"github.com/handmade-jewelry/auth-service/internal/service/service"
	"github.com/handmade-jewelry/auth-service/internal/service/user"
)

type APIHandler struct {
	serviceService  *service.Service
	resourceService *resource.Service
	userService     *user.Service
}

func NewAPIHandler(serviceService *service.Service, resourceService *resource.Service, userService *user.Service) *APIHandler {
	return &APIHandler{
		serviceService:  serviceService,
		resourceService: resourceService,
		userService:     userService,
	}
}
