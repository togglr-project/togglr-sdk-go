package togglr

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Value   string
	Enabled bool
	Found   bool
	Expires time.Time
}

func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.Expires)
}

type LRUCache struct {
	mu       sync.RWMutex
	items    map[string]*CacheEntry
	order    []string
	capacity int
	ttl      time.Duration
}

func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	return &LRUCache{
		items:    make(map[string]*CacheEntry),
		order:    make([]string, 0, capacity),
		capacity: capacity,
		ttl:      ttl,
	}
}

func (c *LRUCache) Get(key string) (*CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.items[key]
	if !exists || entry.IsExpired() {
		return nil, false
	}

	c.moveToEnd(key)

	return entry, true
}

func (c *LRUCache) Set(key string, value string, enabled, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := &CacheEntry{
		Value:   value,
		Enabled: enabled,
		Found:   found,
		Expires: time.Now().Add(c.ttl),
	}

	if _, exists := c.items[key]; exists {
		c.items[key] = entry
		c.moveToEnd(key)

		return
	}

	if len(c.items) >= c.capacity {
		c.evictLRU()
	}

	c.items[key] = entry
	c.order = append(c.order, key)
}

func (c *LRUCache) moveToEnd(key string) {
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)

			break
		}
	}
	c.order = append(c.order, key)
}

func (c *LRUCache) evictLRU() {
	if len(c.order) == 0 {
		return
	}

	key := c.order[0]
	c.order = c.order[1:]
	delete(c.items, key)
}

func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*CacheEntry)
	c.order = c.order[:0]
}

func (c *LRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}
