package s4lru

// from http://www.cs.cornell.edu/~qhuang/papers/sosp_fbanalysis.pdf

import "container/list"

type Item struct {
	lidx  int
	key   string
	value interface{}
}
type Cache struct {
	capacity int
	data     map[string]*list.Element
	lists    []*list.List
}

func New(capacity int) *Cache {
	return &Cache{
		capacity: capacity / 4,
		data:     make(map[string]*list.Element),
		lists:    []*list.List{list.New(), list.New(), list.New(), list.New()},
	}
}

func (c *Cache) Get(key string) interface{} {
	v, ok := c.data[key]

	if !ok {
		return nil
	}

	item := v.Value.(*Item)

	// already on final list?
	if item.lidx == len(c.lists)-1 {
		c.lists[item.lidx].MoveToFront(v)
	} else {
		// move to head of next list
		c.lists[item.lidx].Remove(v)
		delete(c.data, key)
		item.lidx++

		// next list is full, so move the last element of that one to the front of this list
		if c.lists[item.lidx].Len() == c.capacity {
			back := c.lists[item.lidx].Back()
			old := c.lists[item.lidx].Remove(back).(*Item)
			old.lidx--
			c.data[old.key] = c.lists[old.lidx].PushFront(old)
		}

		c.data[key] = c.lists[item.lidx].PushFront(item)
	}

	return item.value
}

func (c *Cache) Put(key string, value interface{}) {
	if c.lists[0].Len() == c.capacity {
		delete(c.data, c.lists[0].Back().Value.(*Item).key)
		c.lists[0].Remove(c.lists[0].Back())
	}
	c.data[key] = c.lists[0].PushFront(&Item{0, key, value})
}
