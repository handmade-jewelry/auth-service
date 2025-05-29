package transport

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/handmade-jewelry/auth-service/internal/service/auth"
	"github.com/handmade-jewelry/auth-service/internal/transport/handler"
	"github.com/handmade-jewelry/auth-service/internal/transport/proxy"
	"github.com/handmade-jewelry/auth-service/logger"
	"github.com/handmade-jewelry/auth-service/pkg/api"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Config struct {
	SwaggerURL          string
	SwaggerSpecURL      string
	SwaggerSpecFilePath string
	HTTPPort            string
}

type Server struct {
	router         chi.Router
	authMiddleware *proxy.AuthMiddleware
	authService    *auth.Service
}

func NewServer(authMiddleware *proxy.AuthMiddleware, authService *auth.Service) *Server {
	return &Server{
		router:         chi.NewRouter(),
		authMiddleware: authMiddleware,
		authService:    authService,
	}
}

func (s *Server) Run(cfg *Config) error {
	s.initSwagger(cfg)

	s.router.Route("/api", func(r chi.Router) {
		r.Use(s.authMiddleware.CheckAccess)

		server := handler.NewAPIHandler(s.authService)
		api.HandlerFromMux(server, r)
	})

	server := handler.NewAPIHandler(s.authService)

	api.HandlerFromMux(server, s.router)

	err := http.ListenAndServe(cfg.HTTPPort, s.router)
	if err != nil {
		logger.Error("error starting server: ", err)
		return err
	}

	logger.Info("HTTP service is running", "port", cfg.HTTPPort)

	return nil
}

func (s *Server) initSwagger(cfg *Config) {
	s.router.HandleFunc(cfg.SwaggerSpecURL, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.SwaggerSpecFilePath)
	})

	s.router.Handle(fmt.Sprintf("%s/*", cfg.SwaggerURL), httpSwagger.Handler(
		httpSwagger.URL(cfg.SwaggerSpecURL),
	))
}
