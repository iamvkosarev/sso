package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/sso/internal/model"
	"github.com/iamvkosarev/sso/internal/storage"
	"github.com/mattn/go-sqlite3"
	"log/slog"
)

type Storage struct {
	db *sql.DB
}

const prepareStorageStatement = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    pass_hash BLOB NOT NULL
);`

// New creates new connection to sqlite3 in provided
// path and returns [Storage] with it.
// If database doesn't exist creates new.
func New(path string, log *slog.Logger) *Storage {
	const op = "storage.sqlite.New"

	log = log.With(slog.String("op", op))

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Error("failed to open database", sl.Err(err))
		return nil
	}
	stmt, err := db.Prepare(prepareStorageStatement)
	if err != nil {
		log.Error("failed to prepare database", sl.Err(err))
		return nil
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Error("failed to prepare database", sl.Err(err))
		return nil
	}
	return &Storage{db: db}
}

func fmtErr(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

func (s Storage) GetUser(email string) (model.User, error) {
	const op = "storage.sqlite.GetUser"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	res, err := stmt.Query(email)
	var sqliteError sqlite3.Error
	if errors.As(err, &sqliteError) {
		if errors.Is(sqliteError.Code, sqlite3.ErrConstraint) {
			return model.User{}, fmt.Errorf("%s: %w", op, storage.ErrorUserNotFound)
		}
		return model.User{}, fmt.Errorf("%s: %w", op, sqliteError)
	}
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer res.Close()
	if res.Next() {
		var user model.User
		if err := res.Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
			return model.User{}, fmt.Errorf("%s: %w", op, err)
		}
		return user, nil
	}
	return model.User{}, storage.ErrorUserNotFound
}

func (s Storage) SaveUser(email string, passHash string) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrorUserExists)
			} else if errors.Is(sqliteErr.Code, sqlite3.ErrConstraintUnique) {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrorUserExists)
			}

			return 0, fmt.Errorf("%s: %w, sql: %w", op, err, sqliteErr)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
