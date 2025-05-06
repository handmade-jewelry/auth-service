package proxy

import (
	"context"
	"encoding/json"
	resourceService "github.com/handmade-jewelry/auth-service/internal/service/resource"
	"net/http"
)

const (
	accessTokenCookie = "access_token"
)

func (a *AuthMiddleware) checkAuth(ctx context.Context, req *http.Request) error {
	// get path
	path := req.URL.Path

	// get handler from redis by path
	resource, err := a.getResource(ctx, path)
	if err != nil {
		//todo error
	}

	// check is_active and deleted_at handler parameters
	if !resource.IsActive && resource.DeletedAt != nil {
		//todo error
	}

	if !resource.CheckAccessToken {
		return nil
	}

	// try to get token from cookies and validate
	token, err := a.getAccessToken(req)
	if err != nil {
		//todo
	}

	// parse token and check sign
	claims, err := a.jwtService.ParseAuthToken(token)
	if err != nil {
		//todo
	}

	if claims.IsExpired() {
		//todo error
	}

	// compare user roles from token with roles from service_handler
	isMatch := a.checkRoles(resource.Roles, claims.Roles)
	if !isMatch {
		//todo error access denied
	}

	return nil
}

func (a *AuthMiddleware) getResource(ctx context.Context, path string) (*resourceService.Resource, error) {
	var resource *resourceService.Resource
	resourceJsn, err := a.redisClient.Get(ctx, path)
	if err == nil {
		err = json.Unmarshal([]byte(resourceJsn), &resource)
		if err == nil {
			return resource, nil
		}
	}

	//todo log redis err

	//try to get from bd
	resource, err = a.resourceService.GetResourceByPath(ctx, path)
	if err != nil {
		return nil, err
	}

	return resource, nil
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
