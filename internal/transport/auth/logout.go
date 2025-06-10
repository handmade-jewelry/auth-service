package auth

import (
	"time"

	"net/http"

	"github.com/handmade-jewelry/auth-service/internal/utils/cookie"
	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
)

func (a *APIHandler) PostLogout(rw http.ResponseWriter, req *http.Request) {
	token, err := cookie.GetCookie(req, cookie.RefreshTokenName)
	if err != nil {
		logger.Error("failed to get cookie refresh token", err)
		errors.WriteHTTPError(rw, errors.UnauthorizedError())
		return
	}

	err = a.authService.Logout(req.Context(), token)
	if err != nil {
		logger.Error("failed to logout", err)
		errors.WriteHTTPError(rw, errors.UnauthorizedError())
		return
	}

	expired := time.Now().Add(-time.Hour)
	http.SetCookie(rw, &http.Cookie{
		Name:     cookie.RefreshTokenName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expired,
		MaxAge:   -1,
	})
	http.SetCookie(rw, &http.Cookie{
		Name:     cookie.AccessTokenName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expired,
		MaxAge:   -1,
	})

	rw.WriteHeader(http.StatusOK)
}
