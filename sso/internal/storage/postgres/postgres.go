package postgres

import (
	"database/sql"
	"fmt"

	"sso/internal/storage"

	"github.com/lib/pq"
)

var (
	RowAlreadyExists = "23505"
)

type Storage struct {
	db *sql.DB
}

func NewPostgresSql(conn string) (*Storage, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) SaveUser(email, password, activationLink string) (int, error) {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (email, password, activation_link) VALUES ($1, $2, $3);")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(email, password, activationLink)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && string(pqErr.Code) == RowAlreadyExists {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserAlreadyExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = s.db.Prepare("SELECT currval('users_id_seq')")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	var userID int
	err = stmt.QueryRow().Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (s *Storage) SaveToken(userID int, token string) error {
	const op = "storage.postgres.SaveToken"

	stmt, err := s.db.Prepare("INSERT INTO tokens (user_id, token) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(userID, token)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
