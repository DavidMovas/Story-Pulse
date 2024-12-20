package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"story-pulse/internal/shared/dbx"
	. "story-pulse/internal/users-service/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByID(ctx context.Context, userId int) (*User, error) {
	builder := dbx.StatementBuilder.
		Select(
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

	query, args, err := builder.ToSql()
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
		return nil, status.Errorf(codes.NotFound, "user %d not found", userId)
	case err != nil:
		return nil, status.Errorf(codes.Internal, "cannot fetch user: %v", err)
	}

	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *UserWithPassword) (*User, error) {
	builder := dbx.StatementBuilder.
		Insert("users").
		Columns("email", "pass_hash", "username").
		Values(user.Email, user.PasswordHash, user.Username).
		Suffix("RETURNING id, role, created_at")

	if user.AvatarURL != nil {
		builder.Columns("avatar_url").Values(*user.AvatarURL)
	}
	if user.FullName != nil {
		builder.Columns("full_name").Values(*user.FullName)
	}
	if user.Bio != nil {
		builder.Columns("bio").Values(*user.Bio)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Role, &user.CreatedAt)

	switch {
	case dbx.IsUniqueViolation(err, "email"):
		return nil, status.Errorf(codes.AlreadyExists, "user with email %s already exists", user.Email)
	case dbx.IsUniqueViolation(err, "username"):
		return nil, status.Errorf(codes.AlreadyExists, "user with username %s already exists", user.Username)
	case err != nil:
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}

	return user.User, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *User) error {
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userId int) error {
	return nil
}
