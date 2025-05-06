package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserID int
	Roles  []string
	jwt.RegisteredClaims
}

type Service struct {
	jwtSecret             []byte
	authTokenExpiryMin    time.Duration
	refreshTokenExpiryMin time.Duration
}

func NewService(jwtSecret string, authTokenExpiryMin, refreshTokenExpiryMin time.Duration) *Service {
	return &Service{
		jwtSecret:             []byte(jwtSecret),
		authTokenExpiryMin:    authTokenExpiryMin,
		refreshTokenExpiryMin: refreshTokenExpiryMin,
	}
}

func (t *Service) GenerateAuthToken(userID int, roles []string) (string, error) {
	claims := Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.authTokenExpiryMin)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *Service) ParseAuthToken(tokenString string) (*Claims, error) {
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

func (c *Claims) IsExpired() bool {
	if c.ExpiresAt == nil {
		return true
	}
	return c.ExpiresAt.Time.Before(time.Now())
}
