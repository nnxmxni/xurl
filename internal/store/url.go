package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type URLMapper struct {
	Url          string `json:"url" validate:"required,url"`
	GeneratedUrl string `json:"generated_url"`
}

type URLStore struct {
	Db *sql.DB
}

func (s *URLStore) Create(ctx context.Context, payload URLMapper) (string, error) {

	query := `
		INSERT INTO url_mapper (short_url, original_url) 
		VALUES ($1, $2) RETURNING short_url
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := s.Db.QueryRowContext(
		ctx,
		query,
		payload.GeneratedUrl,
		payload.Url,
	).Scan(&payload.GeneratedUrl)

	if err != nil {
		return "", err
	}

	return payload.GeneratedUrl, nil
}

func (s *URLStore) Get(ctx context.Context, payload URLMapper) (string, error) {
	query := `
		SELECT original_url FROM url_mapper WHERE short_url = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := s.Db.QueryRowContext(
		ctx,
		query,
		payload.GeneratedUrl,
	).Scan(&payload.Url)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrURLNotFound
		default:
			return "", err
		}
	}

	return payload.Url, nil
}
