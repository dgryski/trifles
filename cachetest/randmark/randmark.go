package randmark

import (
	"math/rand"
)

type element struct {
	value  interface{}
	slot   int
	marked bool
}

type Cache struct {
	capacity int
	data     map[string]element
	marked   []string
	unmarked []string
	mark     bool
}

func New(capacity int) *Cache {
	return &Cache{
		mark:true,
		capacity: capacity,
		data:     make(map[string]element),
		marked:   make([]string, 0, capacity),
		unmarked: make([]string, 0, capacity),
	}
}

func (c *Cache) Get(key string) interface{} {
	e, ok := c.data[key]
	if !ok {
		// no key
		return nil
	}

	if e.marked == !c.mark {
		// key was not marked; we need to mark it

		// remove from unmarked list
		c.removeUnmarked(e.slot)

		// add to marked list
		slot := len(c.marked)
		c.marked = append(c.marked, key)

		// update element
		c.data[key] = element{value: e.value, marked: c.mark, slot: slot}
	}

	return e.value
}

func (c *Cache) removeUnmarked(slot int) {
	last := len(c.unmarked) - 1

	if slot == last {
		c.unmarked = c.unmarked[:last]
		return
	}

	// overwrite slot with last key
	key := c.unmarked[last]
	c.unmarked[slot] = key
	c.unmarked[last] = ""
	c.unmarked = c.unmarked[:last]

	// update slot of moved key
	e := c.data[key]
	e.slot = slot
	c.data[key] = e
}

func (c *Cache) Set(key string, value interface{}) {

	if len(c.data) == c.capacity {
		// need to evict

		// every page is marked
		if len(c.marked) == c.capacity {
			// unmark all the pages
			c.mark = !c.mark
			c.marked, c.unmarked = c.unmarked, c.marked
		}

		// evict unmarked page at random
		slot := rand.Intn(len(c.unmarked))
		delete(c.data, c.unmarked[slot])
		c.removeUnmarked(slot)
	}

	// drop in new key/value as marked
	slot := len(c.marked)
	c.marked = append(c.marked, key)
	c.data[key] = element{value: value, marked: c.mark, slot: slot}
}
