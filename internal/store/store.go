package store

import (
	"context"
	"errors"
)

var ErrURLNotFound = errors.New("resource not found")

type Store interface {
	Create(context.Context, URLMapper) (string, error)
	Get(context.Context, URLMapper) (string, error)
}
