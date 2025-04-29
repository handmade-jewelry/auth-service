package proxy

import (
	rS "github.com/handmade-jewellery/auth-service/internal/service/resource-service"
	sS "github.com/handmade-jewellery/auth-service/internal/service/service-service"
	uS "github.com/handmade-jewellery/auth-service/internal/service/user-service"
	"net/http"
)

type AuthMiddleware struct {
	userService     *uS.UserService
	resourceService *rS.ResourceService
	serviceService  *sS.ServiceService
}

func NewAuthMiddleware(userService *uS.UserService, resourceService *rS.ResourceService,
	serviceService *sS.ServiceService) *AuthMiddleware {
	return &AuthMiddleware{
		userService:     userService,
		resourceService: resourceService,
		serviceService:  serviceService,
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
