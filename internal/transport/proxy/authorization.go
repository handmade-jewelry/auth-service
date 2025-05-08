package proxy

import (
	"context"
	"github.com/pkg/errors"
	"net/http"

	routeService "github.com/handmade-jewelry/auth-service/internal/service/route"
)

const (
	accessTokenCookie = "access_token"
)

func (a *AuthMiddleware) checkAuth(ctx context.Context, req *http.Request) (*routeService.Route, error) {
	path := req.URL.Path

	route, err := a.routeService.GetRouteByPath(ctx, path)
	if err != nil {
		//todo log error
		return nil, err
	}

	if route == nil {
		return nil, errors.New("route not found")
	}

	if !route.CheckAccessToken {
		return route, nil
	}

	token, err := a.getAccessToken(req)
	if err != nil {
		//todo log
		return nil, err
	}

	claims, err := a.jwtService.ParseAuthToken(token)
	if err != nil {
		//todo log
		return nil, err
	}

	if claims.IsExpired() {
		//todo error
		return nil, errors.New("token is expired")
	}

	if !route.CheckRoles {
		return route, nil
	}

	if !a.checkRoles(route.Roles, claims.Roles) {
		//todo error access denied
		return nil, errors.New("access denied")
	}

	return route, nil
}

func (a *AuthMiddleware) getAccessToken(req *http.Request) (string, error) {
	cookie, err := req.Cookie(accessTokenCookie)
	if err != nil {
		//todo
	}
	if len(cookie.Value) == 0 {
		//todo error
	}

	return cookie.Value, nil
}

func (a *AuthMiddleware) checkRoles(resourceRoles, userRoles []string) bool {
	resourceRolesMap := make(map[string]struct{}, len(resourceRoles))
	for _, role := range resourceRoles {
		resourceRolesMap[role] = struct{}{}
	}

	for _, role := range userRoles {
		if _, ok := resourceRolesMap[role]; ok {
			return true
		}
	}

	return false
}
