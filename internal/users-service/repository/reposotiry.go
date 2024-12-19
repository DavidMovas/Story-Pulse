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

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *UserWithPassword) (*User, error) {
	builder := dbx.NewStatementBuilder()

	builder.Insert("users").
		Columns("email", "password_hash", "avatar_url", "username", "full_name", "bio").
		Values(user.Email, user.PasswordHash, user.AvatarURL, user.Username, user.FullName, user.Bio).
		Suffix("RETURNING id, role, created_at")

	query, args, err := builder.Build()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Role, &user.CreatedAt)

	switch {
	case dbx.IsUniqueViolation(err, "email"):
		return nil, apperrors.AlreadyExists("user", "email", user.Email)
	case dbx.IsUniqueViolation(err, "username"):
		return nil, apperrors.AlreadyExists("user", "username", user.Username)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return user.User, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *User) error {
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userId int) error {
	return nil
}
