package repository

import (
	"brain-wave/internal/shared/dbx"
	. "brain-w
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByID(ctx context.Context, userId int) (*User, error) {
	builder := dbx.StatementBuilder.
		Select("id", "email", "avatar_url", "username", "full_name", "bio", "last_login_at", "role", "created_at").
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
		return nil, status.Errorf(codes.NotFound, "user id: %d not found", userId)
	case err != nil:
		return nil, status.Errorf(codes.Internal, "cannot fetch user: %v", err)
	}

	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	builder := dbx.StatementBuilder.
		Select("id, email, avatar_url, username, full_name, bio, pass_hash, last_login_at, role, is_verified, created_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user UserWithPassword
	user.User = &User{}
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.AvatarURL,
		&user.Username,
		&user.FullName,
		&user.Bio,
		&user.PasswordHash,
		&user.LastLoginAt,
		&user.Role,
		&user.IsVerified,
		&user.CreatedAt,
	)

	switch {
	case dbx.IsNoRows(err):
		return nil, status.Errorf(codes.NotFound, "user email: %s not found", email)
	case err != nil:
		return nil, status.Errorf(codes.Internal, "cannot fetch user: %v", err)
	}

	return &user, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*UserWithPassword, error) {
	builder := dbx.StatementBuilder.
		Select("id, email, avatar_url, username, full_name, bio, pass_hash, last_login_at, role, is_verified, created_at").
		From("users").
		Where(squirrel.Eq{"username": username}).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user UserWithPassword
	user.User = &User{}
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.AvatarURL,
		&user.Username,
		&user.FullName,
		&user.Bio,
		&user.PasswordHash,
		&user.LastLoginAt,
		&user.Role,
		&user.IsVerified,
		&user.CreatedAt,
	)

	switch {
	case dbx.IsNoRows(err):
		return nil, status.Errorf(codes.NotFound, "user username: %s not found", username)
	case err != nil:
		return nil, status.Errorf(codes.Internal, "cannot fetch user: %v", err)
	}

	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *UserWithPassword) (*User, error) {
	builder := dbx.StatementBuilder.
		Insert("users").
		Columns("email", "pass_hash", "username").
		Values(user.Email, user.PasswordHash, user.Username).
		Suffix("RETURNING id, role, created_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Role, &user.CreatedAt)

	switch {
	case dbx.IsUniqueViolation(err, "email"):
		return nil, status.Errorf(codes.AlreadyExists, "user email: %s already exists", user.Email)
	case dbx.IsUniqueViolation(err, "username"):
		return nil, status.Errorf(codes.AlreadyExists, "user username: %s already exists", user.Username)
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
