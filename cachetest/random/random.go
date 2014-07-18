package random

// from https://github.com/cloudflare/jgc-talks/tree/master/GoSF/Go_Profiling

import "math/rand"

type Cache struct {
	capacity int
	data     map[string]interface{}
	keys     []string
}

func New(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		data:     make(map[string]interface{}),
		keys:     make([]string, capacity),
	}
}

func (c *Cache) Get(key string) interface{} {
	return c.data[key]
}

func (c *Cache) Set(key string, value interface{}) {

	slot := len(c.data)

	if len(c.data) == c.capacity {

		slot = rand.Intn(c.capacity)

		delete(c.data, c.keys[slot])

	}

	c.keys[slot] = key
	c.data[key] = value
}
