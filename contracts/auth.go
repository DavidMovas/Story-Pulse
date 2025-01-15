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
	User         *User  `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"email"`
	Username string `json:"username" validate:"username"`
	Password string `json:"password" validate:"password"`
}

type LoginUserResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
