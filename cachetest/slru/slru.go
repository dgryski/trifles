package slru

import "github.com/golang/groupcache/lru"

type Cache struct {
	once  *lru.Cache
	twice *lru.Cache
}

func New(oncecap, twicecap int) *Cache {
	return &Cache{
		once:  lru.New(oncecap),
		twice: lru.New(twicecap),
	}
}

func (c *Cache) Get(key string) interface{} {

	val, ok := c.once.Get(key)

	if !ok {
		val, _ = c.twice.Get(key)
		return val
	}

	c.once.Remove(key)
	c.twice.Add(key, val)
	return val
}

func (c *Cache) Set(key string, value interface{}) {
	c.once.Add(key, value)
}
