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

// evictOldest removes the least recently used item from the cache.
// Must be called with write lock held.
func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, item := range c.items {
		if first || item.accessTime.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.accessTime
			first = false
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}

// Delete removes a specific key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*cacheItem)
	c.hits = 0
	c.misses = 0
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Stats returns cache hit/miss statistics.
type Stats struct {
	Hits   int64
	Misses int64
	Ratio  float64
}

// Stats returns the current cache statistics.
func (c *Cache) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var ratio float64
	if c.hits+c.misses > 0 {
		ratio = float64(c.hits) / float64(c.hits+c.misses)
	}

	return Stats{
		Hits:   c.hits,
		Misses: c.misses,
		Ratio:  ratio,
	}
}

// ResetStats resets hit and miss counters to zero.
func (c *Cache) ResetStats() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hits = 0
	c.misses = 0
}

// SetCapacity changes the cache capacity.
// If the new capacity is smaller, evicts oldest items.
func (c *Cache) SetCapacity(capacity int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.capacity = capacity
	for len(c.items) > capacity {
		c.evictOldest()
	}
}

// SetTTL changes the TTL for new items.
func (c *Cache) SetTTL(ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ttl = ttl
}

// Cleanup removes all expired items from the cache.
func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.expiration) {
			delete(c.items, key)
		}
	}
}

// StartCleanupTask starts a background goroutine that periodically cleans up expired items.
func (c *Cache) StartCleanupTask(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			c.Cleanup()
		}
	}()
}
