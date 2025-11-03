package storage

import (
	"context"
	"errors"
)

// خطای دامنه‌ای: پیدا نشد
var ErrNotFound = errors.New("storage: not found")

// تفکیک به interfaceهای کوچک (قابل ترکیب)
type Getter interface {
	Get(ctx context.Context, key string) ([]byte, error)
}
type Putter interface {
	Put(ctx context.Context, key string, val []byte) error
}
type Deleter interface {
	Delete(ctx context.Context, key string) error
}

// Storage = ترکیب سه تک-متدی (کوچک و قابل تست/جایگزینی)
type Storage interface {
	Getter
	Putter
	Deleter
}
