package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) SaveRefreshToken(ctx context.Context, userID int, refreshToken string, expirationTime time.Duration) error {
	key := fmt.Sprintf("refresh_token:%d", userID)

	err := r.client.Set(ctx, key, refreshToken, expirationTime).Err()
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

func (r *Repository) GetRefreshToken(ctx context.Context, userID int) (string, error) {
	key := fmt.Sprintf("refresh_token:%d", userID)

	token, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("no refresh token found for user %d", userID)
	} else if err != nil {
		return "", fmt.Errorf("failed to get refresh token for user %d: %w", userID, err)
	}

	return token, nil
}
