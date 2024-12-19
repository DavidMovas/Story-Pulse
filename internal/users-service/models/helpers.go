package models

func (u *CreateUserRequest) ToUserWithPassword(passwordHash string) *UserWithPassword {
	return &UserWithPassword{
		User: &User{
			Email:     u.Email,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
			FullName:  u.FullName,
			Bio:       u.Bio,
		},
		PasswordHash: passwordHash,
	}
}
