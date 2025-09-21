package togglr

import (
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {
	cache := NewLRUCache(2, 100*time.Millisecond)

	// Test set and get
	cache.Set("key1", "value1", true, true)
	entry, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}
	if entry.Value != "value1" || !entry.Enabled || !entry.Found {
		t.Error("Expected correct values for key1")
	}

	// Test expiration
	time.Sleep(150 * time.Millisecond)
	_, found = cache.Get("key1")
	if found {
		t.Error("Expected key1 to be expired")
	}

	// Test LRU eviction
	cache.Set("key2", "value2", true, true)
	cache.Set("key3", "value3", true, true)
	cache.Set("key4", "value4", true, true) // This should evict key2

	_, found = cache.Get("key2")
	if found {
		t.Error("Expected key2 to be evicted")
	}

	_, found = cache.Get("key3")
	if !found {
		t.Error("Expected key3 to still be present")
	}

	_, found = cache.Get("key4")
	if !found {
		t.Error("Expected key4 to be present")
	}
}

func TestCacheSize(t *testing.T) {
	cache := NewLRUCache(3, time.Hour)

	if cache.Size() != 0 {
		t.Errorf("Expected initial size 0, got %d", cache.Size())
	}

	cache.Set("key1", "value1", true, true)
	if cache.Size() != 1 {
		t.Errorf("Expected size 1, got %d", cache.Size())
	}

	cache.Set("key2", "value2", true, true)
	cache.Set("key3", "value3", true, true)
	if cache.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Size())
	}

	// This should evict key1
	cache.Set("key4", "value4", true, true)
	if cache.Size() != 3 {
		t.Errorf("Expected size 3 after eviction, got %d", cache.Size())
	}
}
