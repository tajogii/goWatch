package cache

import (
	"sync"
	"time"
)

type ICashe[T any] interface {
	Set(key string, value T)
	Get(key string) (T, bool)
}

type Item[T any] struct {
	Value      T
	Expiration int64
}

type Cache[T any] struct {
	items map[string]Item[T]
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewCache[T any](ttl time.Duration) *Cache[T] {
	newCache := &Cache[T]{
		items: make(map[string]Item[T]),
		ttl:   ttl,
	}
	go newCache.cleanUp()
	return newCache

}

func (c *Cache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item[T]{
		Value:      value,
		Expiration: time.Now().Add(c.ttl).UnixNano(),
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()

	item, found := c.items[key]
	if !found {
		var v T
		c.mu.RUnlock()
		return v, false
	}
	c.mu.RUnlock()

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		c.remove(key)
		var v T
		return v, false
	}

	return item.Value, true
}

func (c *Cache[T]) Remove(key string) {
	c.remove(key)
}

func (c *Cache[T]) remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache[T]) cleanUp() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.items {
			if v.Expiration <= time.Now().UnixNano() {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
