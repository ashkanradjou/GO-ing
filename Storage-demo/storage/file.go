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

// نوع concrete unexported
type fileStore struct {
	dir string
}

// سازنده: پوشه را می‌سازد (درصورت نبودن)
func NewFile(dir string) (Storage, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("mkdir %s: %w", dir, err)
	}
	return &fileStore{dir: dir}, nil
}

// کلید را به نام فایل امن تبدیل می‌کنیم (بدون traversal)
func (f *fileStore) keyPath(key string) string {
	// کلید را به base64url بدون padding تبدیل می‌کنیم
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
		// map کردن not-exist به ErrNotFound (و wrap صحیح)
		if errorsIsNotExist(err) {
			// زنجیره شامل ErrNotFound است → errors.Is(err, ErrNotFound) == true
			return nil, fmt.Errorf("file get %q: %w", key, ErrNotFound)
		}
		// سایر خطاهای I/O: wrap خود خطای اصلی
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

	// نوشتن اتمی: به فایل موقت بنویس، بعد rename
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

// errorsIsNotExist: کمک‌کننده برای تشخیص not-exist با errors.As روی *fs.PathError
func errorsIsNotExist(err error) bool {
	// راه کوتاه:
	// return os.IsNotExist(err)

	// برای نمایش errors.As:
	var pe *fs.PathError
	if ok := errorsAs(err, &pe); ok {
		return os.IsNotExist(pe)
	}
	return os.IsNotExist(err)
}

// لایه‌ی نازک قابل جایگزینی برای تست/دمو (تا در import حلقه ایجاد نشود)
var (
	errorsAs = func(err error, target any) bool { return errors.As(err, target) }
)
