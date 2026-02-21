package cache

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c := New(10, time.Minute)
	if c.capacity != 10 {
		t.Errorf("expected capacity 10, got %d", c.capacity)
	}
	if c.ttl != time.Minute {
		t.Errorf("expected ttl 1m, got %v", c.ttl)
	}
}

func TestCache_SetAndGet(t *testing.T) {
	c := New(10, time.Minute)
	
	// Test Set
	c.Set("key1", "value1")
	
	// Test Get
	val, found := c.Get("key1")
	if !found {
		t.Error("expected to find key1")
	}
	if val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}
}

func TestCache_GetNotFound(t *testing.T) {
	c := New(10, time.Minute)
	
	val, found := c.Get("nonexistent")
	if found {
		t.Error("expected not to find nonexistent key")
	}
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}
}

func TestCache_Delete(t *testing.T) {
	c := New(10, time.Minute)
	
	c.Set("key1", "value1")
	c.Delete("key1")
	
	val, found := c.Get("key1")
	if found {
		t.Error("expected key1 to be deleted")
	}
	if val != nil {
		t.Errorf("expected nil after delete, got %v", val)
	}
}

func TestCache_Clear(t *testing.T) {
	c := New(10, time.Minute)
	
	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Clear()
	
	if c.Len() != 0 {
		t.Errorf("expected length 0 after clear, got %d", c.Len())
	}
}

func TestCache_Len(t *testing.T) {
	c := New(10, time.Minute)
	
	c.Set("key1", "value1")
	c.Set("key2", "value2")
	
	if c.Len() != 2 {
		t.Errorf("expected length 2, got %d", c.Len())
	}
}

func TestCache_Expiration(t *testing.T) {
	c := New(10, time.Millisecond)
	
	c.Set("key1", "value1")
	
	// Wait for expiration
	time.Sleep(10 * time.Millisecond)
	
	val, found := c.Get("key1")
	if found {
		t.Error("expected key1 to be expired")
	}
	if val != nil {
		t.Errorf("expected nil after expiration, got %v", val)
	}
}

func TestCache_EvictOldest(t *testing.T) {
	c := New(2, time.Minute)
	
	c.Set("key1", "value1")
	c.Set("key2", "value2")
	// This should evict key1
	c.Set("key3", "value3")
	
	// key1 should be evicted
	_, found := c.Get("key1")
	if found {
		t.Error("expected key1 to be evicted")
	}
	
	// key2 and key3 should exist
	_, found = c.Get("key2")
	if !found {
		t.Error("expected key2 to exist")
	}
	
	_, found = c.Get("key3")
	if !found {
		t.Error("expected key3 to exist")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := New(100, time.Minute)
	
	// Run concurrent operations
	done := make(chan bool)
	
	// Writer
	go func() {
		for i := 0; i < 100; i++ {
			c.Set(string(rune(i)), i)
		}
		done <- true
	}()
	
	// Reader
	go func() {
		for i := 0; i < 100; i++ {
			c.Get(string(rune(i)))
		}
		done <- true
	}()
	
	<-done
	<-done
}

func TestCache_EmptyDelete(t *testing.T) {
	c := New(10, time.Minute)
	
	// Delete from empty cache should not panic
	c.Delete("nonexistent")
}

func TestCache_EmptyLen(t *testing.T) {
	c := New(10, time.Minute)
	
	if c.Len() != 0 {
		t.Errorf("expected length 0, got %d", c.Len())
	}
}
