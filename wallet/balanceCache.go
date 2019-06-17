package wallet

import (
	"sync"
	"time"
)

// BalanceCache holds temporary data as a buffer
//   ttl - the duration a value should be cached
type BalanceCache struct {
	ttl       time.Duration
	mutex     *sync.RWMutex
	cache     map[string]uint64
	cacheTime map[string]time.Time
}

// NewBalanceCache create a new default cache with a ttl of 60 seconds
func NewBalanceCache() *BalanceCache {
	return &BalanceCache{
		60 * time.Second,
		&sync.RWMutex{},
		make(map[string]uint64),
		make(map[string]time.Time),
	}
}

// Clear removes all data from the cache
func (c *BalanceCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = make(map[string]uint64)
	c.cacheTime = make(map[string]time.Time)
}

// Get retrieve a value for the given key, returns false in second
// param if key not found or expired
func (c *BalanceCache) Get(key string) (uint64, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if time.Since(c.cacheTime[key]) <= c.ttl {
		return c.cache[key], true
	}
	return 0, false
}

// Set store a value for the given key
func (c *BalanceCache) Set(key string, val uint64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[key] = val
	c.cacheTime[key] = time.Now()
}
