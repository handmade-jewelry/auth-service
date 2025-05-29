package proxy

import (
	"github.com/handmade-jewelry/auth-service/internal/jwt"
	routeService "github.com/handmade-jewelry/auth-service/internal/service/route"
	userService "github.com/handmade-jewelry/auth-service/internal/service/user"
	"net/http"
	"net/http/httputil"
)

type AuthMiddleware struct {
	userService  *userService.Service
	routeService *routeService.Service
	jwtService   *jwt.Service
}

func NewAuthMiddleware(
	userService *userService.Service,
	routeService *routeService.Service,
	jwtService *jwt.Service,
) *AuthMiddleware {
	return &AuthMiddleware{
		userService:  userService,
		routeService: routeService,
		jwtService:   jwtService,
	}
}

func (a *AuthMiddleware) CheckAccess(_ http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		route, err := a.checkAuth(ctx, req)
		if err != nil {
			//todo log and process error
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

		proxy.ServeHTTP(wr, req)
	})
}
