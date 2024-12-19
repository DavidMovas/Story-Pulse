package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"story-pulse/internal/shared/dbx"
	apperrors "story-pulse/internal/shared/error"
	. "story-pulse/internal/users-service/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByID(ctx context.Context, userId int) (*User, error) {
	builder := dbx.NewStatementBuilder()
	builder.Select(
		"id",
		"email",
		"avatar_url",
		"username",
		"full_name",
		"bio",
		"last_login_at",
		"role",
		"created_at",
	).
		From("users").
		Where(squirrel.Eq{"id": userId}).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.Build()

	if err != nil {
		return nil, err
	}

	var user User
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.AvatarURL,
		&user.Username,
		&user.FullName,
		&user.Bio,
		&user.LastLoginAt,
		&user.Role,
		&user.CreatedAt,
	)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "id", userId)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return &user, nil
}
