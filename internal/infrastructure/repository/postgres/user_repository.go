package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/iamvkosarev/sso/internal/domain/entity"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: pool}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	const op = "repository.postgres.GetByEmail"

	query := `SELECT id, email, pass_hash FROM users WHERE email = $1`

	var user entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.Id, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) Save(ctx context.Context, user entity.User) (entity.UserId, error) {
	const op = "repository.postgres.Save"

	var id int64
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id`,
		user.Email, user.PassHash,
	).Scan(&id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, entity.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, pgErr)
	}

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return entity.UserId(id), nil
}
