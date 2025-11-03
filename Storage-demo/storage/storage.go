package storage

import (
	"context"
	"errors"
)

// Domain error: not found
var ErrNotFound = errors.New("storage: not found")

// Separation into small interfaces (combinable)
type Getter interface {
	Get(ctx context.Context, key string) ([]byte, error)
}
type Putter interface {
	Put(ctx context.Context, key string, val []byte) error
}
type Deleter interface {
	Delete(ctx context.Context, key string) error
}

// Storage = Combination of three single-methods (small and testable/replaceable)
type Storage interface {
	Getter
	Putter
	Deleter
}
