package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/handmade-jewelry/auth-service/internal/cache"
	"github.com/handmade-jewelry/auth-service/internal/jwt"
	"github.com/handmade-jewelry/auth-service/internal/service/user"
)

const refreshTokenPrefix = "refresh_token:"

type Service struct {
	jwtService      *jwt.Service
	redisClient     *cache.RedisClient
	userService     *user.Service
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewService(
	jwtService *jwt.Service,
	redisClient *cache.RedisClient,
	userService *user.Service,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Service {
	return &Service{
		jwtService:      jwtService,
		redisClient:     redisClient,
		userService:     userService,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *Service) RefreshToken(ctx context.Context, token string) (*AuthTokens, error) {
	claims, err := s.jwtService.ParseRefreshToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	oldTokenID := claims.ID
	userID := claims.UserID
	err = s.validateRefreshToken(ctx, oldTokenID, userID)
	if err != nil {
		return nil, fmt.Errorf("refresh token is not valid: %w", err)
	}

	roles, err := s.userService.UserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	tokenID := uuid.NewString()
	authTokens, err := s.generateAuthTokens(tokenID, userID, roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth tokens: %w", err)
	}

	err = s.tokenRotate(ctx, oldTokenID, tokenID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to rotate refresh tokens: %w", err)
	}

	return authTokens, nil
}

func (s *Service) validateRefreshToken(ctx context.Context, oldTokenID string, tokenUserID int64) error {
	val, err := s.redisClient.Get(ctx, refreshTokenPrefix+oldTokenID)
	if err == redis.Nil {
		return fmt.Errorf("invalid or expired refresh token: %w", err)
	} else if err != nil {
		return fmt.Errorf("invalid redis get: %w", err)
	}

	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fmt.Errorf("userID convert errors: %w", err)
	}

	if tokenUserID != userID {
		return fmt.Errorf("token user ID mismatch: token=%d redis=%d", tokenUserID, userID)
	}

	return nil
}

func (s *Service) generateAuthTokens(tokenID string, userID int64, roles []string) (*AuthTokens, error) {
	accessToken, err := s.jwtService.GenerateAccessToken(userID, roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(userID, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RefreshTTL:   s.refreshTokenTTL,
		AccessTTL:    s.accessTokenTTL,
	}, nil
}

func (s *Service) tokenRotate(ctx context.Context, oldTokenID, newTokenID string, userID int64) error {
	keyOld := refreshTokenPrefix + oldTokenID
	keyNew := refreshTokenPrefix + newTokenID
	userIDStr := strconv.FormatInt(userID, 10)

	pipe := s.redisClient.Pipeline()
	pipe.Del(ctx, keyOld)
	pipe.Set(ctx, keyNew, userIDStr, s.refreshTokenTTL)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis pipeline exec failed: %w", err)
	}

	for _, cmd := range cmds {
		if cmd.Err() != nil && cmd.Err() != redis.Nil {
			return fmt.Errorf("redis command failed: %w", cmd.Err())
		}
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*AuthTokens, error) {
	userWithRoles, err := s.userService.Login(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("user login failed: %w", err)
	}

	tokenID := uuid.NewString()
	authTokens, err := s.generateAuthTokens(tokenID, userWithRoles.UserID, userWithRoles.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth tokens: %w", err)
	}

	err = s.redisClient.Set(ctx, refreshTokenPrefix+tokenID, userWithRoles.UserID, s.refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token in redis: %w", err)
	}

	return authTokens, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.jwtService.ParseRefreshToken(refreshToken)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	err = s.redisClient.Delete(ctx, refreshTokenPrefix+claims.ID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token from redis: %w", err)
	}

	return nil
}
