package proxy

import (
	"github.com/handmade-jewelry/auth-service/internal/cache"
	"github.com/handmade-jewelry/auth-service/internal/jwt"
	resourceService "github.com/handmade-jewelry/auth-service/internal/service/resource"
	serviceService "github.com/handmade-jewelry/auth-service/internal/service/service"
	userService "github.com/handmade-jewelry/auth-service/internal/service/user"
	"net/http"
)

type AuthMiddleware struct {
	userService     *userService.Service
	resourceService *resourceService.Service
	serviceService  *serviceService.Service
	redisClient     *cache.RedisClient
	jwtService      *jwt.Service
}

func NewAuthMiddleware(
	userService *userService.Service,
	resourceService *resourceService.Service,
	serviceService *serviceService.Service,
	jwtService *jwt.Service,
	redisClient *cache.RedisClient,
) *AuthMiddleware {
	return &AuthMiddleware{
		userService:     userService,
		resourceService: resourceService,
		serviceService:  serviceService,
		redisClient:     redisClient,
		jwtService:      jwtService,
	}
}

func (a *AuthMiddleware) CheckAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		err := a.checkAuth(ctx, req)
		if err != nil {
			//todo log and process error
			test := err
			_ = test
		}

		//todo is it better to send a new HTTP request to filter only the data you need
		next.ServeHTTP(wr, req)
	})
}
