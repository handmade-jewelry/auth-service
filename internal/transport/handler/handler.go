package handler

import (
	userService "github.com/handmade-jewelry/auth-service/internal/service/user"
)

// todo rename
type Handler struct {
	userService *userService.Service
}

func NewHandler() *Handler {
	return &Handler{}
}
