package models

type RefreshToken struct {
	UserID int    `json:"user_id" db:"user_id"`
	Role   string `json:"role" db:"role"`
}
