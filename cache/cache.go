package cache

import (
	"container/list"
	"errors"
	"fmt"
)

type keyT uint32
type valueT string

type elem struct {
	key keyT
	val valueT
}

type Cache struct {
	Capacity int
	table    map[keyT]*list.Element
	queue    *list.List
}

func NewCache(cap int) (*Cache, error) {
	if cap <= 0 {
		return nil, errors.New("capacity must be a natural number")
	}
	c := &Cache{
		Capacity: cap,
		table:    make(map[keyT]*list.Element),
		queue:    list.New(),
	}
	return c, nil
}

func (c *Cache) Put(key keyT, val valueT) {
	if element, ok := c.table[key]; ok {
		c.queue.MoveToFront(element)
		element.Value.(*elem).val = val
		return
	}

	c.table[key] = c.queue.PushFront(&elem{key, val})

	if c.queue.Len() > c.Capacity {
		c.RemoveElem(c.queue.Back())
	}
}

func (c *Cache) RemoveElem(element *list.Element) {
	delete(c.table, element.Value.(*elem).key)
	c.queue.Remove(element)
}

func (c *Cache) Get(key keyT) (val valueT, ok bool) {
	if element, ok := c.table[key]; ok {
		c.queue.MoveToFront(element)
		return element.Value.(*elem).val, true
	}
	return "", false
}

func (c *Cache) Desplay() {
	str := "{"
	for e := c.queue.Front(); e != nil; e = e.Next() {
		str += fmt.Sprintf("{%v: %v}, ", e.Value.(*elem).key, e.Value.(*elem).val)
	}
	str = str[:len(str)-2]
	str += "}"
	fmt.Println(str)
}
