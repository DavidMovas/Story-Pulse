package models

import (
	grpc "brain-wave/internal/shared/grpc/v1"
	"brain-wave/internal/shared/helpers"
)

func (u *User) ToGRPC() *grpc.User {
	return &grpc.User{
		Id:          int64(u.ID),
		Email:       u.Email,
		AvatarUrl:   u.AvatarURL,
		Username:    u.Username,
		FullName:    u.FullName,
		Bio:         u.Bio,
		Role:        u.Role,
		LastLoginAt: helpers.ToTimestamp(u.LastLoginAt),
		CreatedAt:   helpers.ToTimestamp(&u.CreatedAt),
	}
}

func ToUserWithPassword(r *grpc.CreateUserRequest) *UserWithPassword {
	return &UserWithPassword{
		User: &User{
			Email:    r.Email,
			Username: r.Username,
		},
		PasswordHash: r.Password,
	}
}
