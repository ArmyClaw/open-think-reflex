package cache

import (
	"sync"
	"time"
)

// Cache is a simple LRU cache
type Cache struct {
	mu       sync.RWMutex
	items    map[string]*cacheItem
	capacity int
	ttl      time.Duration
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// New creates a new cache
func New(capacity int, ttl time.Duration) *Cache {
	return &Cache{
		items:    make(map[string]*cacheItem),
		capacity: capacity,
		ttl:      ttl,
	}
}

// Get gets a value from cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check expiration
	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Set sets a value in cache
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict if at capacity
	if len(c.items) >= c.capacity {
		c.evictOldest()
	}

	c.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}

// Delete deletes a value from cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear clears the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*cacheItem)
}

// evictOldest removes the oldest item from cache
func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range c.items {
		if oldestTime.IsZero() || item.expiration.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.expiration
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}

// Len returns the number of items in cache
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
