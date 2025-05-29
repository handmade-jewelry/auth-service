package auth

import (
	"github.com/handmade-jewelry/auth-service/internal/utils/cookie"
	"github.com/handmade-jewelry/auth-service/logger"
	"github.com/handmade-jewelry/auth-service/pkg/api/auth"
	"net/http"
	"time"
)

func (a *APIHandler) GetRefreshToken(wr http.ResponseWriter, req *http.Request, _ auth.GetRefreshTokenParams) {
	token, err := cookie.GetCookie(req, cookie.RefreshTokenName)
	if err != nil {
		logger.Error("failed to get refresh token cookie", err)
		http.Error(wr, "refresh token not found", http.StatusUnauthorized)
		return
	}

	authTokens, err := a.authService.RefreshToken(req.Context(), token)
	if err != nil {
		logger.Error("failed to logout user", err)
		http.Error(wr, "unauthorized", http.StatusUnauthorized)
		return
	}

	now := time.Now()
	http.SetCookie(wr, &http.Cookie{
		Name:     cookie.RefreshTokenName,
		Value:    authTokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  now.Add(authTokens.RefreshTTL),
	})

	http.SetCookie(wr, &http.Cookie{
		Name:     cookie.AccessTokenName,
		Value:    authTokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  now.Add(authTokens.AccessTTL),
	})

	wr.WriteHeader(http.StatusOK)
}
