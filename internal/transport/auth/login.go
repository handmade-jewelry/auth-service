package auth

import (
	"encoding/json"
	"github.com/handmade-jewelry/auth-service/internal/utils/cookie"
	"github.com/handmade-jewelry/auth-service/logger"
	"net/http"
	"time"
)

func (a *APIHandler) PostLogin(wr http.ResponseWriter, req *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var data loginRequest
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		http.Error(wr, "Invalid request body", http.StatusBadRequest)
		return
	}

	if data.Email == "" || data.Password == "" {
		http.Error(wr, "Invalid email or password", http.StatusBadRequest)
		return
	}

	authTokens, err := a.authService.Login(req.Context(), data.Email, data.Password)
	if err != nil {
		logger.ErrorWithFields(
			"failed to login user",
			err,
			"email", data.Email,
			"password", data.Password)
		http.Error(wr, "Unauthorized", http.StatusUnauthorized)
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
