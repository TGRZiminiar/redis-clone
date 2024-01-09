package cache

import "sync"

type (
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

// func (c *Cache)
