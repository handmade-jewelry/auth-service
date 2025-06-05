package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/handmade-jewelry/auth-service/errors"
	"github.com/handmade-jewelry/auth-service/internal/jwt"
	routeService "github.com/handmade-jewelry/auth-service/internal/service/route"
)

type AuthMiddleware struct {
	routeService *routeService.Service
	jwtService   *jwt.Service
}

func NewAuthMiddleware(
	routeService *routeService.Service,
	jwtService *jwt.Service,
) *AuthMiddleware {
	return &AuthMiddleware{
		routeService: routeService,
		jwtService:   jwtService,
	}
}

func (a *AuthMiddleware) CheckAccess(_ http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		route, err := a.checkAuth(ctx, req)
		if err != nil {
			errors.WriteHTTPError(rw, err)
			return
		}

		// todo check proxy
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = route.Scheme
				req.URL.Host = route.Host
				req.Host = route.Host
				req.URL.Path = route.ServicePath
			},
		}

		proxy.ServeHTTP(rw, req)
	})
}
