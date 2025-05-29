package transport

import (
	"github.com/go-chi/chi/v5"
	"github.com/handmade-jewelry/auth-service/internal/transport/auth"
	"github.com/handmade-jewelry/auth-service/internal/transport/proxy"
	"github.com/handmade-jewelry/auth-service/internal/transport/resource"
	"github.com/handmade-jewelry/auth-service/logger"
	pkgAuth "github.com/handmade-jewelry/auth-service/pkg/api/auth"
	pkgGateway "github.com/handmade-jewelry/auth-service/pkg/api/resource"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type SwaggerConfig struct {
	SwaggerURL              string
	SwaggerAuthURL          string
	SwaggerResourceURL      string
	SwaggerAuthSpecURL      string
	SwaggerAuthSpecPath     string
	SwaggerResourceSpecURL  string
	SwaggerResourceSpecPath string
}

type Opts struct {
	HTTPPort       string
	ProxyPrefix    string
	AuthPrefix     string
	ResourcePrefix string
}

type Server struct {
	opts              *Opts
	router            chi.Router
	authMiddleware    *proxy.AuthMiddleware
	authAPIHandler    *auth.APIHandler
	gatewayAPIHandler *resource.APIHandler
}

func NewServer(
	opts *Opts,
	authMiddleware *proxy.AuthMiddleware,
	authAPIHandler *auth.APIHandler,
	gatewayAPIHandler *resource.APIHandler,
) *Server {
	return &Server{
		opts:              opts,
		router:            chi.NewRouter(),
		authMiddleware:    authMiddleware,
		authAPIHandler:    authAPIHandler,
		gatewayAPIHandler: gatewayAPIHandler,
	}
}

func (s *Server) Run(cfg *SwaggerConfig) error {
	s.initSwagger(cfg)

	s.router.Route(s.opts.ProxyPrefix, func(r chi.Router) {
		r.Route(s.opts.AuthPrefix, func(r chi.Router) {
			pkgAuth.HandlerFromMux(s.authAPIHandler, r)
		})

		r.Route(s.opts.ResourcePrefix, func(r chi.Router) {
			pkgGateway.HandlerFromMux(s.gatewayAPIHandler, r)
		})

		r.Group(func(r chi.Router) {
			r.Use(s.authMiddleware.CheckAccess)
			r.NotFound(http.NotFoundHandler().ServeHTTP)
			r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			})
		})
	})

	err := http.ListenAndServe(s.opts.HTTPPort, s.router)
	if err != nil {
		logger.Error("error starting server: ", err)
		return err
	}

	logger.Info("HTTP service is running", "port", s.opts.HTTPPort)

	return nil
}

func (s *Server) initSwagger(cfg *SwaggerConfig) {
	s.router.HandleFunc(cfg.SwaggerAuthSpecURL, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.SwaggerAuthSpecPath)
	})

	s.router.HandleFunc(cfg.SwaggerResourceSpecURL, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.SwaggerResourceSpecPath)
	})

	s.router.Handle(cfg.SwaggerAuthURL, httpSwagger.Handler(
		httpSwagger.URL(cfg.SwaggerAuthSpecURL),
	))

	s.router.Handle(cfg.SwaggerResourceURL, httpSwagger.Handler(
		httpSwagger.URL(cfg.SwaggerResourceSpecURL),
	))
}
