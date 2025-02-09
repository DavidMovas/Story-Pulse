package contracts

import "time"

type User struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	AvatarURL   *string    `json:"avatarUrl,omitempty"`
	Username    string     `json:"username"`
	FullName    *string    `json:"fullName,omitempty"`
	Bio         *string    `json:"bio,omitempty"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	Role        string     `json:"role"`
	IsVerified  *bool      `json:"isVerified,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type UserShort struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatarUrl,omitempty"`
	FullName  *string   `json:"fullName,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type UserWithPassword struct {
	*User
	PasswordHash string `json:"passwordHash"`
}

type UserDetails struct {
	*User
	ArticlesCount  *int `json:"articlesCount,omitempty"`
	CommentsCount  *int `json:"commentsCount,omitempty"`
	FollowersCount *int `json:"followersCount,omitempty"`
	FollowingCount *int `json:"followingCount,omitempty"`
}

type GetUserByIDRequest struct {
	ID string `json:"-" param:"userId" validate:"required"`
}
