// Package cache provides a small in-memory TTL cache with single-flight
// semantics. Underlying data only changes when a models.dev sync completes,
// so expensive aggregate endpoints (stats, providers, sitemap) are cached
// and invalidated on successful sync.
package cache

import (
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

type Cache struct {
	mu       sync.Mutex
	entries  map[string]*entry
	inFlight map[string]chan struct{}
}

func New() *Cache {
	return &Cache{
		entries:  make(map[string]*entry),
		inFlight: make(map[string]chan struct{}),
	}
}

// GetOrCompute returns the cached value for key if present and unexpired.
// Otherwise it invokes fn; concurrent callers for the same key share a
// single fn invocation. Results are only cached when fn returns no error.
func (c *Cache) GetOrCompute(key string, ttl time.Duration, fn func() (any, error)) (any, error) {
	c.mu.Lock()
	if e, ok := c.entries[key]; ok && time.Now().Before(e.expiresAt) {
		v := e.value
		c.mu.Unlock()
		return v, nil
	}
	if done, ok := c.inFlight[key]; ok {
		c.mu.Unlock()
		<-done
		return c.GetOrCompute(key, ttl, fn)
	}
	done := make(chan struct{})
	c.inFlight[key] = done
	c.mu.Unlock()

	value, err := fn()

	c.mu.Lock()
	delete(c.inFlight, key)
	if err == nil && value != nil {
		c.entries[key] = &entry{value: value, expiresAt: time.Now().Add(ttl)}
	}
	close(done)
	c.mu.Unlock()

	return value, err
}

// Invalidate drops all cached values. Called after a successful data sync.
func (c *Cache) Invalidate() {
	c.mu.Lock()
	c.entries = make(map[string]*entry)
	c.mu.Unlock()
}
