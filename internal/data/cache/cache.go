package cache

import (
	"sync"
	"time"
)

// Cache is a simple LRU cache with statistics tracking
type Cache struct {
	mu       sync.RWMutex
	items    map[string]*cacheItem
	capacity int
	ttl      time.Duration
	
	// Statistics
	hits   int64
	misses int64
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
	accessTime time.Time
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
		c.misses++
		return nil, false
	}

	// Check expiration
	if time.Now().After(item.expiration) {
		c.misses++
		return nil, false
	}

	// Update access time for LRU
	item.accessTime = time.Now()
	c.hits++
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
		accessTime: time.Now(),
	}
}

// GetOrSet gets a value from cache or computes and stores it
func (c *Cache) GetOrSet(key string, compute func() interface{}) interface{} {
	// Try to get from cache first (read lock)
	c.mu.RLock()
	item, found := c.items[key]
	if found && !time.Now().After(item.expiration) {
		c.mu.RUnlock()
		c.mu.Lock()
		c.hits++
		c.mu.Unlock()
		return item.value
	}
	c.mu.RUnlock()
	
	// Compute the value
	value := compute()
	
	// Store in cache
	c.Set(key, value)
	
	return value
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
	c.hits = 0
	c.misses = 0
}

// evictOldest removes the oldest item from cache
func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range c.items {
		if oldestTime.IsZero() || item.accessTime.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.accessTime
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

// Stats returns cache statistics
func (c *Cache) Stats() (hits, misses int64, ratio float64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	hits = c.hits
	misses = c.misses
	
	total := hits + misses
	if total > 0 {
		ratio = float64(hits) / float64(total)
	}
	
	return
}

// ResetStats resets the cache statistics
func (c *Cache) ResetStats() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hits = 0
	c.misses = 0
}

// SetCapacity sets a new capacity and evicts items if necessary
func (c *Cache) SetCapacity(capacity int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.capacity = capacity
	
	// Evict if over capacity
	for len(c.items) > capacity {
		c.evictOldest()
	}
}

// SetTTL sets a new TTL for new items
func (c *Cache) SetTTL(ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ttl = ttl
}

// Cleanup removes expired items from the cache
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

// StartCleanupTask starts a background task to periodically clean up expired items
func (c *Cache) StartCleanupTask(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			c.Cleanup()
		}
	}()
}
