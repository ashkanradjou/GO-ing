package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ashkanradjou/GO-ing/storage-demo/storage"
)

func runDemo(ctx context.Context, st storage.Storage) error {
	key := "user:42"
	val := []byte("hello")

	// Put
	if err := st.Put(ctx, key, val); err != nil {
		return fmt.Errorf("put failed: %w", err)
	}

	// Get
	got, err := st.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("get failed: %w", err)
	}
	fmt.Println("GET:", string(got))

	// Delete
	if err := st.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	// Get after delete → should give ErrNotFound
	_, err = st.Get(ctx, key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			fmt.Println("GET after delete → not found (as expected)")
		} else {
			return fmt.Errorf("unexpected get error: %w", err)
		}
	}
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Memory
	fmt.Println("== Memory Store ==")
	mem := storage.NewMemory()
	if err := runDemo(ctx, mem); err != nil {
		fmt.Println("ERROR:", err)
	}

	// File
	fmt.Println("\n== File Store ==")
	fs, err := storage.NewFile("./data")
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	if err := runDemo(ctx, fs); err != nil {
		fmt.Println("ERROR:", err)
	}
}
