package models

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	grpc "story-pulse/internal/shared/grpc/v1"
)

func (u *User) ToGRPC() *grpc.User {
	return &grpc.User{
		Id:          int64(u.ID),
		Email:       u.Email,
		AvatarUrl:   *u.AvatarURL,
		Username:    u.Username,
		FullName:    *u.FullName,
		Bio:         *u.Bio,
		LastLoginAt: timestamppb.New(*u.LastLoginAt),
		Role:        u.Role,
		CreatedAt:   timestamppb.New(u.CreatedAt),
	}
}

func ToUserWithPassword(r *grpc.CreateUserRequest) *UserWithPassword {
	return &UserWithPassword{
		User: &User{
			Email:     r.Email,
			Username:  r.Username,
			AvatarURL: &r.AvatarUrl,
			FullName:  &r.FullName,
			Bio:       &r.Bio,
		},
		PasswordHash: r.Password,
	}
}
