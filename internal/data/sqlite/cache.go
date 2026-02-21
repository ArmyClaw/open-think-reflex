package sqlite

import (
	"sync"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// QueryCache caches frequent database queries for performance
type QueryCache struct {
	mu           sync.RWMutex
	patternCache map[string]*cacheEntry
	stats        CacheStats
	ttl          time.Duration
	maxSize     int
}

type cacheEntry struct {
	pattern   *models.Pattern
	expiresAt time.Time
}

type CacheStats struct {
	Hits   int64
	Misses int64
	Evicts int64
}

// NewQueryCache creates a new query cache
func NewQueryCache(ttl time.Duration, maxSize int) *QueryCache {
	return &QueryCache{
		patternCache: make(map[string]*cacheEntry),
		ttl:          ttl,
		maxSize:     maxSize,
	}
}

// Get retrieves a cached pattern by ID
func (c *QueryCache) Get(id string) (*models.Pattern, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.patternCache[id]
	if !exists || time.Now().After(entry.expiresAt) {
		c.stats.Misses++
		return nil, false
	}

	c.stats.Hits++
	return entry.pattern, true
}

// Set stores a pattern in cache
func (c *QueryCache) Set(id string, pattern *models.Pattern) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict if at capacity
	if len(c.patternCache) >= c.maxSize {
		c.evictOldest()
	}

	c.patternCache[id] = &cacheEntry{
		pattern:   pattern,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Invalidate removes a pattern from cache
func (c *QueryCache) Invalidate(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.patternCache, id)
}

// Clear empties the cache
func (c *QueryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.patternCache = make(map[string]*cacheEntry)
}

func (c *QueryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.patternCache {
		if oldestTime.IsZero() || entry.expiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.expiresAt
		}
	}

	if oldestKey != "" {
		delete(c.patternCache, oldestKey)
		c.stats.Evicts++
	}
}

// Stats returns cache statistics
func (c *QueryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats
}
