package tworand

// "power of two random choices" -- pick 2 random elements, remove the oldest

import "math/rand"

type entry struct {
	key string
	val interface{}
	t   int64
}

type Cache struct {
	capacity int
	keys     map[string]int
	data     []entry
	t        int64
}

func New(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		keys:     make(map[string]int),
		data:     make([]entry, capacity),
	}
}

func (c *Cache) Get(key string) interface{} {
	if v, ok := c.keys[key]; ok {
		c.data[v].t = c.t
		return c.data[v].val
	}

	return nil
}

func (c *Cache) Set(key string, value interface{}) {

	slot := len(c.data)

	if len(c.data) == c.capacity {

		slot = rand.Intn(c.capacity)
		if slot2 := rand.Intn(c.capacity); c.data[slot2].t < c.data[slot].t {
			slot = slot2
		}

		delete(c.keys, c.data[slot].key)
	}

	c.keys[key] = slot
	c.data[slot].key = key
	c.data[slot].val = value
	c.data[slot].t = c.t
	c.t++
}
