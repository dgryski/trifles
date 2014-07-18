package clock

type entry struct {
	key   string
	value interface{}
	used  bool
}

type Cache struct {
	capacity int
	data     map[string]int
	keys     []entry
	hand     int
}

func New(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		data:     make(map[string]int),
		keys:     make([]entry, capacity),
	}
}

func (c *Cache) Get(key string) interface{} {
	idx, ok := c.data[key]
	if !ok {
		return nil
	}
	c.keys[idx].used = true
	return c.keys[idx].value
}

func (c *Cache) Set(key string, value interface{}) {

	slot := len(c.data)

	if slot == c.capacity {

		for c.keys[c.hand].used {
			c.keys[c.hand].used = false
			c.hand++
			if c.hand == c.capacity {
				c.hand = 0
			}
		}

		slot = c.hand

		delete(c.data, c.keys[slot].key)

	}

	c.keys[slot].key = key
	c.keys[slot].value = value
	c.keys[slot].used = true
	c.data[key] = slot
}
