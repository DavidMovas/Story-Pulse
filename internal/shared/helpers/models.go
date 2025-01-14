package helpers

import (
	"story-pulse/contracts"
	grpc "story-pulse/internal/shared/grpc/v1"
	"strconv"
)

func ToGRPC(u contracts.User) *grpc.User {
	id, _ := strconv.Atoi(u.ID)
	return &grpc.User{
		Id:          int64(id),
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
