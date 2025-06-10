package auth

import (
	"time"

	"encoding/json"
	"net/http"

	"github.com/handmade-jewelry/auth-service/internal/utils/cookie"
	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
)

func (a *APIHandler) PostLogin(rw http.ResponseWriter, req *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var data loginRequest
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		errors.WriteHTTPError(rw, errors.BadRequestError("Invalid request body"))
		return
	}

	if data.Email == "" || data.Password == "" {
		http.Error(rw, "Invalid email or password", http.StatusBadRequest)
		return
	}

	authTokens, err := a.authService.Login(req.Context(), data.Email, data.Password)
	if err != nil {
		logger.ErrorWithFields(
			"failed to login user",
			err,
			"email", data.Email,
			"password", data.Password)

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
