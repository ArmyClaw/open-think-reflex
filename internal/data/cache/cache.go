// Package cache provides an in-memory LRU (Least Recently Used) cache
// with TTL (Time To Live) support and statistics tracking.
// Thread-safe using RWMutex for concurrent read/write access.
package cache

import (
	"sync"
	"time"
)

// Cache is a thread-safe LRU cache with TTL expiration.
// Implements Get/Set operations with automatic eviction of expired
// or least recently used items.
type Cache struct {
	mu       sync.RWMutex
	items    map[string]*cacheItem
	capacity int           // Maximum number of items
	ttl      time.Duration // Time-to-live for each item

	// Statistics - track cache performance
	hits   int64 // Number of successful cache hits
	misses int64 // Number of cache misses
}

// cacheItem stores a cached value with metadata
type cacheItem struct {
	value      interface{}    // The cached value
	expiration time.Time     // When this item expires
	accessTime time.Time     // Last access time (for LRU eviction)
}

// New creates a new cache with the specified capacity and TTL.
//   - capacity: maximum number of items ( eviction occurs when exceeded)
//   - ttl: time-to-live for each item (0 = no expiration)
func New(capacity int, ttl time.Duration) *Cache {
	return &Cache{
		items:    make(map[string]*cacheItem),
		capacity: capacity,
		ttl:      ttl,
	}
}

// Get retrieves a value from the cache.
// Returns (value, true) if found and not expired.
// Returns (nil, false) if not found or expired.
// Updates hit/miss statistics.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		c.misses++
		return nil, false
	}

	// Check expiration
	if time.Now().After(item.expiration) {
		c.misses++
		return nil, false
	}

	// Update access time for LRU tracking
	item.accessTime = time.Now()
	c.hits++
	return item.value, true
}

// Set stores a value in the cache with TTL.
// If capacity is exceeded, evicts the least recently used item.
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
		accessTime: time.Now(),
	}
}

// GetOrSet retrieves a value from cache, or computes and stores it if missing.
// The compute function is called only if the item is not in cache or expired.
// This is atomic - only one goroutine will compute for a given key.
func (c *Cache) GetOrSet(key string, compute func() interface{}) interface{} {
	// Try to get from cache first (read lock)
	c.mu.RLock()
	item, found := c.items[key]
	if found && !time.Now().After(item.expiration) {
		item.accessTime = time.Now()
		c.hits++
		c.mu.RUnlock()
		return item.value
	}
	c.mu.RUnlock()

	// Compute the value (needs write lock)
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	item, found = c.items[key]
	if found && !time.Now().After(item.expiration) {
		item.accessTime = time.Now()
		c.hits++
		return item.value
	}

	// Compute and store
	value := compute()
	c.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl),
		accessTime: time.Now(),
	}

	// Evict if necessary
	if len(c.items) > c.capacity {
		c.evictOldest()
	}

	return value
}
