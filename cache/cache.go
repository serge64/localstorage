package cache

import "sync/atomic"

type Cache struct {
	a atomic.Value
}

// New creates a new instance Cache.
func New() Cache {
	return Cache{atomic.Value{}}
}

// the Save signals that the cache
// has been saved.
func (c *Cache) Save() {
	c.a.Store(true)
}

// the Reset signals that the cache
// has been reset.
func (c *Cache) Reset() {
	c.a.Store(false)
}

// Cached returns true if the cache
// exists.
func (c *Cache) Cached() bool {
	cache := c.a.Load()
	if cache == nil {
		return false
	}
	return cache.(bool)
}
