package cache

import (
	"container/list"
	"errors"
	"fmt"
)

type KeyT uint32
type ValueT string

type elem struct {
	key KeyT
	val ValueT
}

type Cache struct {
	сapacity uint16 // range 1..29'999
	table    map[KeyT]*list.Element
	queue    *list.List
}

// Creates a new Cache with a fixed size in range: 1..29'999
func NewCache(cap uint16) (*Cache, error) {
	if cap <= 0 || cap >= 30000 {
		return nil, errors.New("capacity must be a natural number, less than 30000")
	}
	c := &Cache{
		сapacity: cap,
		table:    make(map[KeyT]*list.Element),
		queue:    list.New(),
	}
	return c, nil
}

// Puts element: {key, value} to Cache
func (c *Cache) Put(key KeyT, val ValueT) {
	if element, ok := c.table[key]; ok {
		c.queue.MoveToFront(element)
		element.Value.(*elem).val = val
		return
	}

	c.table[key] = c.queue.PushFront(&elem{key, val})

	if uint16(c.queue.Len()) > c.сapacity {
		c.removeElem(c.queue.Back())
	}
}

// Returns the value by key, true - if the item exists.
// If the key is not in the Cache - nil, false.
func (c *Cache) Get(key KeyT) (val ValueT, ok bool) {
	if element, ok := c.table[key]; ok {
		c.queue.MoveToFront(element)
		return element.Value.(*elem).val, true
	}
	return "", false
}

// Displays the contents of the cache.
func (c *Cache) Display() {
	str := "{"
	for e := c.queue.Front(); e != nil; e = e.Next() {
		str += fmt.Sprintf("{%v: %v}, ", e.Value.(*elem).key, e.Value.(*elem).val)
	}
	str = str[:len(str)-2]
	str += "}"
	fmt.Println(str)
}

func (c *Cache) removeElem(element *list.Element) {
	delete(c.table, element.Value.(*elem).key)
	c.queue.Remove(element)
}
