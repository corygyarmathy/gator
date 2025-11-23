// Package cachegator provides a cache for package api.
package cachegator

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte // Raw cached data
}

type Cache struct {
	cacheEntries map[string]cacheEntry
	mu           sync.Mutex
	done         chan struct{}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cacheEntries: make(map[string]cacheEntry),
		done:         make(chan struct{}),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheEntries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.cacheEntries[key]
	if !ok {
		return nil, ok
	}
	return e.val, ok
}

func (c *Cache) Close() {
	select {
	case <-c.done:
		// done is already closed; return without closing again
		return
	default:
		// done is not closed; safe to close it now
		close(c.done)
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cutoff := time.Now().UTC().Add(-interval)
			c.reap(cutoff)
		case <-c.done:
			return
		}
	}
}

func (c *Cache) reap(cutoff time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.cacheEntries {
		if v.createdAt.Before(cutoff) {
			delete(c.cacheEntries, k)
		}
	}
}
