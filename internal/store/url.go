package store

import "database/sql"

type URLStore struct {
	Db *sql.DB
}

func (s *URLStore) Create(url, slug string) (string, error) {
	return "", nil
}
