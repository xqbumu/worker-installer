package handler

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type TestResult struct {
	Value string
}

func TestInMemoryCache(t *testing.T) {
	cache := NewInMemoryCache[string](2)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if val, ok := cache.Get("key1"); !ok || val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}

	cache.Set("key3", "value3")

	if _, ok := cache.Get("key2"); ok {
		t.Errorf("expected key2 to be evicted")
	}
}

func TestFileCache(t *testing.T) {
	dir := os.TempDir()
	cache, err := NewFileCache[string](dir, time.Second, 2, 1)
	if err != nil {
		t.Fatalf("failed to create file cache: %v", err)
	}

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if val, ok := cache.Get("key1"); !ok || val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}

	cache.Set("key3", "value3")

	if _, ok := cache.Get("key1"); !ok {
		t.Errorf("expected key1 to be present")
	}

	cache.CleanUp(time.Second)
	time.Sleep(time.Second * 1)

	if _, ok := cache.Get("key2"); ok {
		t.Errorf("expected key2 to be evicted")
	}

	// Test file persistence
	cache.Set("key4", "value4")
	cache.flushBatch()

	filePath := filepath.Join(dir, "key4.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("expected key4 to be saved to file")
	}

	// Clean up
	os.Remove(filePath)
}
