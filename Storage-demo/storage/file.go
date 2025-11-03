package storage

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// concrete unexported
type fileStore struct {
	dir string
}

// Creator: Creates the folder (if it doesn't exist)
func NewFile(dir string) (Storage, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("mkdir %s: %w", dir, err)
	}
	return &fileStore{dir: dir}, nil
}

// Convert the key to a secure file name (no traversal)
func (f *fileStore) keyPath(key string) string {
	//Convert the key to base64url without padding
	name := base64.RawURLEncoding.EncodeToString([]byte(key))
	return filepath.Join(f.dir, name+".bin")
}

func (f *fileStore) Get(ctx context.Context, key string) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	p := f.keyPath(key)
	b, err := os.ReadFile(p)
	if err != nil {
		//Map not-exist to ErrNotFound (and wrap correctly)
		if errorsIsNotExist(err) {
			//The chain contains ErrNotFound → errors.Is(err, ErrNotFound) == true
			return nil, fmt.Errorf("file get %q: %w", key, ErrNotFound)
		}
		// Other I/O errors: wrap the main error itself
		return nil, fmt.Errorf("file get %q: %w", key, err)
	}
	return b, nil
}

func (f *fileStore) Put(ctx context.Context, key string, val []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	p := f.keyPath(key)

	// Atomic write: write to temporary file, then rename
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, val, 0o600); err != nil {
		return fmt.Errorf("file put write tmp %q: %w", key, err)
	}
	if err := os.Rename(tmp, p); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("file put rename %q: %w", key, err)
	}
	return nil
}

func (f *fileStore) Delete(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	p := f.keyPath(key)
	if err := os.Remove(p); err != nil {
		if errorsIsNotExist(err) {
			return fmt.Errorf("file delete %q: %w", key, ErrNotFound)
		}
		return fmt.Errorf("file delete %q: %w", key, err)
	}
	return nil
}

// errorsIsNotExist: Helper for detecting not-exist with errors.As on *fs.PathError
func errorsIsNotExist(err error) bool {

	var pe *fs.PathError
	if ok := errorsAs(err, &pe); ok {
		return os.IsNotExist(pe)
	}
	return os.IsNotExist(err)
}

// Replaceable thin layer for testing/demo (to avoid looping in import)
var (
	errorsAs = func(err error, target any) bool { return errors.As(err, target) }
)
