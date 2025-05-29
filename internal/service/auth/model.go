package auth

import "time"

type AuthTokens struct {
	RefreshToken string
	AccessToken  string
	RefreshTTL   time.Duration
	AccessTTL    time.Duration
}
