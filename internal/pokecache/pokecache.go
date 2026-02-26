package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	// use a mutex to protect access to a shared variable (ie prevent race conditions)
	mtx sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{}
	c.entries = make(map[string]cacheEntry)
	go c.reapLoop(interval)
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.entries[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {

	c.mtx.Lock()
	defer c.mtx.Unlock()

	val, ok := c.entries[key]
	if ok {
		return val.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	// ticker.C i think just means this for loop runs every time the ticker ticks. double check this though.
	for range ticker.C {
		// check each entry in the cache and remove any that have expired
		c.mtx.Lock()
		for key, entry := range c.entries {
			age := time.Since(entry.createdAt)
			if age > interval {
				// we reap
				delete(c.entries, key)
			}
		}
		c.mtx.Unlock()
	}
}
