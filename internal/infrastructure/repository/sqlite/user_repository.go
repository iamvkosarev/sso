package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/iamvkosarev/sso/internal/domain/entity"
	"github.com/mattn/go-sqlite3"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	const op = "repository.sqlite.GetUser"

	stmt, err := r.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	res, err := stmt.Query(email)
	var sqliteError sqlite3.Error
	if errors.As(err, &sqliteError) {
		if errors.Is(sqliteError.Code, sqlite3.ErrConstraint) {
			return entity.User{}, fmt.Errorf("%s: %w", op, entity.ErrUserExists)
		}
		return entity.User{}, fmt.Errorf("%s: %w", op, sqliteError)
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer res.Close()
	if res.Next() {
		var user entity.User
		if err := res.Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
			return entity.User{}, fmt.Errorf("%s: %w", op, err)
		}
		return user, nil
	}
	return entity.User{}, entity.ErrUserNotFound
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) Save(ctx context.Context, user entity.User) (int64, error) {
	result, err := r.db.ExecContext(
		ctx, "INSERT INTO users (email, pass_hash) VALUES (?, ?)", user.Email, user.PassHash,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
