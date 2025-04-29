package transport

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/handmade-jewellery/auth-service/internal/transport/handler"
	"github.com/handmade-jewellery/auth-service/internal/transport/proxy"
	"github.com/handmade-jewellery/auth-service/pkg/api"
	"log"
	"net/http"
)

type Server struct {
	router         chi.Router
	authMiddleware *proxy.AuthMiddleware
}

func NewServer(authMiddleware *proxy.AuthMiddleware) *Server {
	return &Server{
		router:         chi.NewRouter(),
		authMiddleware: authMiddleware,
	}
}

func (s *Server) Run() error {
	fmt.Println("Proxy service is running on :8080")

	s.router.Use(s.authMiddleware.CheckAccess)

	server := handler.NewHandler()
	api.HandlerFromMux(server, s.router)

	//todo addr from config.yaml
	err := http.ListenAndServe(":8090", s.router)
	if err != nil {
		log.Printf("Error starting server: %v\n", err)
		return err
	}

	return nil
}
