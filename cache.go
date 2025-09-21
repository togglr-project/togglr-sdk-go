package togglr

import (
	"sync"
	"time"
)

// CacheEntry represents a cached evaluation result
type CacheEntry struct {
	Value   string
	Enabled bool
	Found   bool
	Expires time.Time
}

// IsExpired checks if the cache entry has expired
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.Expires)
}

// LRUCache implements a simple LRU cache with TTL
type LRUCache struct {
	mu       sync.RWMutex
	items    map[string]*CacheEntry
	order    []string
	capacity int
	ttl      time.Duration
}

// NewLRUCache creates a new LRU cache
func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	return &LRUCache{
		items:    make(map[string]*CacheEntry),
		order:    make([]string, 0, capacity),
		capacity: capacity,
		ttl:      ttl,
	}
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(key string) (*CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.items[key]
	if !exists || entry.IsExpired() {
		return nil, false
	}

	// Move to end (most recently used)
	c.moveToEnd(key)

	return entry, true
}

// Set stores a value in the cache
func (c *LRUCache) Set(key string, value string, enabled, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := &CacheEntry{
		Value:   value,
		Enabled: enabled,
		Found:   found,
		Expires: time.Now().Add(c.ttl),
	}

	// If a key exists, update and move to end
	if _, exists := c.items[key]; exists {
		c.items[key] = entry
		c.moveToEnd(key)

		return
	}

	// If at capacity, remove least recently used
	if len(c.items) >= c.capacity {
		c.evictLRU()
	}

	// Add a new entry
	c.items[key] = entry
	c.order = append(c.order, key)
}

// moveToEnd moves the key to the end of the order slice
func (c *LRUCache) moveToEnd(key string) {
	// Find and remove from the current position
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)

			break
		}
	}
	// Add to end
	c.order = append(c.order, key)
}

// evictLRU removes the least recently used item
func (c *LRUCache) evictLRU() {
	if len(c.order) == 0 {
		return
	}

	// Remove the first item (least recently used)
	key := c.order[0]
	c.order = c.order[1:]
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*CacheEntry)
	c.order = c.order[:0]
}

// Size returns the current number of items in the cache
func (c *LRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}
