package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Service struct {
	jwtSecret          []byte
	authTokenExpiry    time.Duration
	refreshTokenExpiry time.Duration
}

func NewService(jwtSecret string, authTokenExpiry, refreshTokenExpiry time.Duration) *Service {
	return &Service{
		jwtSecret:          []byte(jwtSecret),
		authTokenExpiry:    authTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

func (t *Service) GenerateAccessToken(userID int64, roles []string) (string, error) {
	claims := ClaimsWithRoles{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.authTokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *Service) ParseAccessToken(tokenString string) (*ClaimsWithRoles, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClaimsWithRoles{}, func(token *jwt.Token) (interface{}, error) {
		return t.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ClaimsWithRoles)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

func (t *Service) GenerateRefreshToken(userID int64, tokenID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.refreshTokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *Service) ParseRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return t.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
