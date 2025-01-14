package contracts

type RegisterUserRequest struct {
	Email     string  `json:"email" validate:"email"`
	Password  string  `json:"password" validate:"password"`
	Username  string  `json:"username" validate:"username"`
	AvatarURL *string `json:"avatarUrl"`
	FullName  *string `json:"fullName"`
	Bio       *string `json:"bio"`
}

type RegisterUserResponse struct {
	User         *User  `json:"user,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
