package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type ClaimsWithRoles struct {
	UserID int64
	Roles  []string
	jwt.RegisteredClaims
}

type Claims struct {
	UserID int64
	jwt.RegisteredClaims
}

func (c *ClaimsWithRoles) IsExpired() bool {
	if c.ExpiresAt == nil {
		return true
	}
	return c.ExpiresAt.Time.Before(time.Now())
}
