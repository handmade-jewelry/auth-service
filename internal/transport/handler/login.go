package handler

import (
	"github.com/handmade-jewelry/auth-service/internal/utils/cookie"
	"net/http"
	"time"
)

func (a *APIHandler) PostLogin(wr http.ResponseWriter, req *http.Request) {
	authTokens, err := a.authService.Login(req.Context())
	if err != nil {
		//todo
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
