package repository

import (
	"brain-wave/internal/auth-service/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	refreshTokenPattern = "refresh_token:%s"
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) SaveRefreshToken(ctx context.Context, userID int, role, refreshToken string, expirationTime time.Duration) error {
	key := fmt.Sprintf(refreshTokenPattern, refreshToken)

	data := &models.RefreshToken{
		UserID: userID,
		Role:   role,
	}

	jsonData, err := json.Marshal(data)

	err = r.client.Set(ctx, key, jsonData, expirationTime).Err()
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

func (r *Repository) ValidateRefreshToken(ctx context.Context, refreshToken string) (*models.RefreshToken, error) {
	key := fmt.Sprintf(refreshTokenPattern, refreshToken)

	data, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("invalid or expired refresh token")
	} else if err != nil {
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	var token models.RefreshToken
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal refresh token: %w", err)
	}

	return &token, nil
}

func (r *Repository) RemoveRefreshToken(ctx context.Context, refreshToken string) error {
	key := fmt.Sprintf(refreshTokenPattern, refreshToken)
	return r.client.Del(ctx, key).Err()
}
