package postgres

import (
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func NewPostgresSql(conn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
