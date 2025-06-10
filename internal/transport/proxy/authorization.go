package proxy

import (
	"context"

	"net/http"

	routeService "github.com/handmade-jewelry/auth-service/internal/service/route"
	"github.com/handmade-jewelry/auth-service/internal/utils/cookie"
	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
)

func (a *AuthMiddleware) checkAuth(ctx context.Context, req *http.Request) (*routeService.Route, *errors.HTTPError) {
	path := req.URL.Path

	route, err := a.routeService.GetRouteByPath(ctx, path)
	if err != nil {
		logger.ErrorWithFields("failed to get route", err, "path", path)
		return nil, errors.NotFoundError("Resource not found")
	}

	if !route.CheckAccessToken {
		return route, nil
	}

	token, err := cookie.GetCookie(req, cookie.AccessTokenName)
	if err != nil {
		logger.ErrorWithFields("failed to get access token from cookie", err, "path", path)
		return nil, errors.UnauthorizedError()
	}

	claims, err := a.jwtService.ParseAccessToken(token)
	if err != nil {
		logger.Error("failed to parse token", err)
		return nil, errors.UnauthorizedError()
	}

	if claims.IsExpired() {
		return nil, errors.UnauthorizedError()
	}

	if !route.CheckRoles {
		return route, nil
	}

	if !a.checkRoles(route.Roles, claims.Roles) {
		logger.Error("roles is mismatched", err)
		return nil, errors.UnauthorizedError()
	}

	return route, nil
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
