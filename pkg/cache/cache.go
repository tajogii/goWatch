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
	items      map[string]Item[T]
	mu         sync.RWMutex
	expiration int64
}

func NewCache[T any](ttl time.Duration) *Cache[T] {
	newCache := &Cache[T]{
		items:      make(map[string]Item[T]),
		expiration: time.Now().Add(ttl).UnixNano(),
	}
	return newCache

}

func (c *Cache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item[T]{
		Value:      value,
		Expiration: c.expiration,
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		var v T
		return v, false
	}

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		c.Remove(key)
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
