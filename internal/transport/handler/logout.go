package handler

import (
	"github.com/handmade-jewelry/auth-service/internal/utils/cookie"
	"github.com/handmade-jewelry/auth-service/logger"
	"net/http"
	"time"
)

func (a *APIHandler) PostLogout(wr http.ResponseWriter, req *http.Request) {
	token, err := cookie.GetCookie(req, cookie.RefreshTokenName)
	if err != nil {
		logger.Error("failed to get refresh token cookie", err)
		http.Error(wr, "refresh token not found", http.StatusUnauthorized)
		return
	}

	err = a.authService.Logout(req.Context(), token)
	if err != nil {
		http.Error(wr, "", http.StatusUnauthorized)
		return
	}

	expired := time.Now().Add(-time.Hour)
	http.SetCookie(wr, &http.Cookie{
		Name:     cookie.RefreshTokenName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expired,
		MaxAge:   -1,
	})
	http.SetCookie(wr, &http.Cookie{
		Name:     cookie.AccessTokenName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expired,
		MaxAge:   -1,
	})

	wr.WriteHeader(http.StatusOK)
}
