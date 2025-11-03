package storage

import (
	"context"
	"fmt"
	"sync"
)

// concrete unexported
type memStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

// Exported constructor; concrete type does not leak
func NewMemory() Storage {
	return &memStore{
		data: make(map[string][]byte),
	}
}

func (m *memStore) Get(ctx context.Context, key string) ([]byte, error) {
	// Respect context cancellation (fast here)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.data[key]
	if !ok {
		// wrap with ErrNotFound for errors.Is
		return nil, fmt.Errorf("mem get %q: %w", key, ErrNotFound)
	}
	//Copy for safety
	out := make([]byte, len(v))
	copy(out, v)
	return out, nil
}

func (m *memStore) Put(ctx context.Context, key string, val []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Copy for safety
	buf := make([]byte, len(val))
	copy(buf, val)
	m.data[key] = buf
	return nil
}

func (m *memStore) Delete(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; !ok {
		return fmt.Errorf("mem delete %q: %w", key, ErrNotFound)
	}
	delete(m.data, key)
	return nil
}
