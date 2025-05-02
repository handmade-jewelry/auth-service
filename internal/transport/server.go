package transport

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/handmade-jewelry/auth-service/internal/transport/handler"
	"github.com/handmade-jewelry/auth-service/internal/transport/proxy"
	"github.com/handmade-jewelry/auth-service/pkg/api"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

type Config struct {
	SwaggerURLPath      string
	SwaggerSpecFilePath string
	HTTPPort            string
}

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

func (s *Server) Run(cfg *Config) error {
	s.router.Use(s.authMiddleware.CheckAccess)

	s.initSwagger(cfg)

	server := handler.NewHandler()

	api.HandlerFromMux(server, s.router)

	err := http.ListenAndServe(cfg.HTTPPort, s.router)
	if err != nil {
		log.Printf("Error starting server: %v\n", err)
		return err
	}

	fmt.Println("Proxy service is running on :8090")

	return nil
}

func (s *Server) initSwagger(cfg *Config) {
	s.router.Get(cfg.SwaggerURLPath, httpSwagger.Handler(
		httpSwagger.URL(cfg.SwaggerURLPath),
	))

	s.router.Get(cfg.SwaggerURLPath, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.SwaggerURLPath)
	})
}
