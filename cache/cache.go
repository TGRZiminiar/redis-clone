package cache

import (
	"fmt"
	"sync"
	"time"
)

type (
	Cacher interface {
		Set([]byte, []byte, time.Duration) error
		Has([]byte) bool
		Get([]byte) ([]byte, error)
		Delete([]byte) error
	}

	Cache struct {
		// Use RWmutex cause most of the time we read instead of write
		mu   sync.RWMutex
		data map[string][]byte
	}
)

func NewCache() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Delete(key []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, string(key))
	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.data[string(key)]
	return ok
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keyStr := string(key)
	val, ok := c.data[keyStr]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", keyStr)
	}

	return val, nil
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[string(key)] = value

	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			delete(c.data, string(key))
		}()
	}

	return nil
}
