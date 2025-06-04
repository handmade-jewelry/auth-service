package proxy

import (
	"context"
	"fmt"
	"github.com/handmade-jewelry/auth-service/internal/utils/cookie"
	"github.com/handmade-jewelry/auth-service/logger"
	"net/http"

	routeService "github.com/handmade-jewelry/auth-service/internal/service/route"
)

// todo errors
func (a *AuthMiddleware) checkAuth(ctx context.Context, req *http.Request) (*routeService.Route, error) {
	path := req.URL.Path

	route, err := a.routeService.GetRouteByPath(ctx, path)
	if err != nil {
		logger.ErrorWithFields("failed to get route", err, "path", path)
		return nil, fmt.Errorf("failed to get route: %w", err)
	}

	if !route.CheckAccessToken {
		return route, nil
	}

	token, err := cookie.GetCookie(req, cookie.AccessTokenName)
	if err != nil {
		logger.Error("failed to get access token from cookie", err)
		return nil, fmt.Errorf("failed to get cookie: %w", err)
	}

	claims, err := a.jwtService.ParseAccessToken(token)
	if err != nil {
		logger.Error("failed to parse token", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims.IsExpired() {
		logger.Error("token is expired", err)
		return nil, fmt.Errorf("token is expired: %w", err)
	}
	//todo надо ли обнулить куки в таком кейсе?

	//todo обновлять access token прямо в middleware, как только он протух
	//drop refresh-token handler

	if !route.CheckRoles {
		return route, nil
	}

	//todo
	//сверять список ролей с закешированным чтобы проверить что все роли валидные

	if !a.checkRoles(route.Roles, claims.Roles) {
		logger.Error("roles is mismatched", err)
		return nil, fmt.Errorf("access denied: %w", err)
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
