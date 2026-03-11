package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.pg.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: failed to ping DB: %w", op, err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS url(
        id SERIAL PRIMARY KEY,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL
    );
`)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create table: %w", op, err)
	}

	_, err = db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_alias 
    ON url(alias);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create index: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.pg.SaveURL"

	var id int64
	err := s.db.QueryRow("INSERT INTO url(url, alias) VALUES ($1, $2) RETURNING id", urlToSave, alias).Scan(&id)
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.pg.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)

	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.pg.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: alias %s not found", op, alias)
	}

	return nil
}
