package proxy

import (
	"context"
	"net/http"
)

const (
	authTokenHeader = "auth_token"
)

func (a *AuthMiddleware) checkAuth(ctx context.Context, req *http.Request) error {
	path := req.URL.Path

	//todo get resource from cache
	resource, err := a.resourceService.GetRouteByPath(ctx, path)
	if err != nil {
		return err
	}
	//todo validate resource - active/deleted

	//todo process case resource without roles

	cookie, err := req.Cookie(authTokenHeader)
	if err != nil {
		//todo
	}
	//todo validate cookie token

	userRoles, err := a.userService.GetUserRoles(ctx, cookie.Value)
	if err != nil {
		//todo
	}

	err = a.checkRoles(resource.Roles, userRoles.Roles)
	if err != nil {
		//todo
	}

	return nil
}

func (a *AuthMiddleware) checkRoles(resourceRoles, userRoles []string) error {
	//todo stub
	return nil
}
