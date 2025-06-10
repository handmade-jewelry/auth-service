package auth

import (
	"time"

	"net/http"

	"github.com/handmade-jewelry/auth-service/internal/utils/cookie"
	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
	"github.com/handmade-jewelry/auth-service/pkg/api/auth"
)

func (a *APIHandler) GetRefreshToken(rw http.ResponseWriter, req *http.Request, _ auth.GetRefreshTokenParams) {
	token, err := cookie.GetCookie(req, cookie.RefreshTokenName)
	if err != nil {
		logger.Error("failed to get cookie refresh token", err)
		errors.WriteHTTPError(rw, errors.UnauthorizedError())
		return
	}

	authTokens, err := a.authService.RefreshToken(req.Context(), token)
	if err != nil {
		logger.Error("failed to refresh token", err)
		errors.WriteHTTPError(rw, errors.UnauthorizedError())
		return
	}

	now := time.Now()
	http.SetCookie(rw, &http.Cookie{
		Name:     cookie.RefreshTokenName,
		Value:    authTokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  now.Add(authTokens.RefreshTTL),
	})

	http.SetCookie(rw, &http.Cookie{
		Name:     cookie.AccessTokenName,
		Value:    authTokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  now.Add(authTokens.AccessTTL),
	})

	rw.WriteHeader(http.StatusOK)
}
