package helpers

import (
	"story-pulse/contracts"
	grpc "story-pulse/internal/shared/grpc/v1"
)

func ToGRPC(u contracts.User) *grpc.User {
	return &grpc.User{
		Id:          int64(u.ID),
		Email:       u.Email,
		AvatarUrl:   u.AvatarURL,
		Username:    u.Username,
		FullName:    u.FullName,
		Bio:         u.Bio,
		Role:        u.Role,
		LastLoginAt: ToTimestamp(u.LastLoginAt),
		CreatedAt:   ToTimestamp(&u.CreatedAt),
	}
}
