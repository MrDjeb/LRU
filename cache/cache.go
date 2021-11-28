package cache

import (
	"container/list"
	"errors"
	"fmt"
	"math"
	. "reflect"
	"runtime"
	"sync"
	"time"
)

const (
	DefaultTTL     = 30              //time.Second
	DefaultMinSize = 1 << 20         // 1Mb
	DefaultMaxSize = (1 << 20) * 128 // 128Mb
)

type KeyT uint32
type ValueT string

type elem struct {
	key      KeyT
	val      ValueT
	deadTime int64
}

func (e *elem) Size() uint32 {
	size := (uint32(ValueOf(e.key).Type().Size()) + 16) * 2
	size += uint32(len(e.val))
	size += 32
	size += uint32(ValueOf(e.deadTime).Type().Size())
	return size
}

type Cache struct {
	table         map[KeyT]*list.Element
	queue         *list.List
	ttl           uint32
	checkTTL      uint32
	curSize       uint32
	minSize       uint32
	maxSize       uint32
	isCollapse    chan bool
	isClose       chan bool
	closeCollapse sync.WaitGroup
	sync.RWMutex
}

// Creates a new Cache.
// Input free or less argument: TTL, minByteSize, maxByteSize
func NewCache(param ...uint32) (*Cache, error) {
	if len(param) > 3 {
		return nil, errors.New("number of arguments must be less than 3")
	}
	param = append(param, DefaultTTL, DefaultMinSize, DefaultMaxSize)
	fl := 0
	if len(param) == 1 {
		fl = 1
	}
	c := &Cache{
		table:         make(map[KeyT]*list.Element),
		queue:         list.New(),
		ttl:           param[0],
		checkTTL:      uint32(math.Sqrt(float64(param[0]))),
		curSize:       0,
		minSize:       param[1+fl],
		maxSize:       param[2+len(param)%3],
		isCollapse:    make(chan bool, 1),
		isClose:       make(chan bool),
		closeCollapse: sync.WaitGroup{},
		RWMutex:       sync.RWMutex{},
	}
	c.closeCollapse.Add(1)
	go c.collapse()
	return c, nil
}

// Puts element: {key, value} to Cache
func (c *Cache) Put(key KeyT, val ValueT) {
	c.Lock()
	defer c.Unlock()

	if e, ok := c.table[key]; ok {
		c.queue.MoveToFront(e)
		c.curSize -= e.Value.(*elem).Size()
		e.Value.(*elem).val = val
		e.Value.(*elem).deadTime = time.Now().Add(time.Duration(c.ttl) * time.Second).Unix()
		c.curSize += e.Value.(*elem).Size()
	} else {
		c.table[key] = c.queue.PushFront(&elem{key, val,
			time.Now().Add(time.Duration(c.ttl) * time.Second).Unix()})
		c.curSize += c.queue.Front().Value.(*elem).Size()
		if c.curSize >= c.maxSize {
			c.isCollapse <- true
		}
	}
}

// Returns the value by key, true - if the item exists.
// If the key is not in the Cache - nil, false.
func (c *Cache) Get(key KeyT) (val ValueT, ok bool) {
	c.Lock()
	defer c.Unlock()

	if e, ok := c.table[key]; ok {
		if time.Now().UnixNano() > e.Value.(*elem).deadTime {
			return "", false
		}
		c.queue.MoveToFront(e)
		e.Value.(*elem).deadTime = time.Now().Add(time.Duration(c.ttl) * time.Second).Unix()
		return e.Value.(*elem).val, true
	}
	return "", false
}

// Displays the contents of the cache.
func (c *Cache) Display() {
	c.RLock()
	defer c.RUnlock()

	str := "{"
	for e := c.queue.Front(); e != nil; e = e.Next() {
		str += fmt.Sprintf("{%v: %v}, ", e.Value.(*elem).key, e.Value.(*elem).val)
	}
	if str != "{" {
		str = str[:len(str)-2]
	}
	str += "}"
	fmt.Println(str)
}

func (c *Cache) Destroy() {
	if c != nil {
		c.Lock()
		if c.isClose != nil {
			close(c.isClose)
			c.closeCollapse.Wait()
			c.isClose = nil
		}
		c.Unlock()
	}
	c = nil
}

func (c *Cache) removeElem(e *list.Element) {
	c.Lock()
	c.curSize -= e.Value.(*elem).Size()
	delete(c.table, e.Value.(*elem).key)
	c.queue.Remove(e)
	c.Unlock()
}

func (c *Cache) collapse() {
	defer c.closeCollapse.Done()
	timer := time.NewTicker(time.Second * time.Duration(c.checkTTL))
	defer timer.Stop()

	for {
		select {
		case <-c.isCollapse:
			for e := c.queue.Back(); c.curSize <= c.maxSize && e != nil; e = e.Prev() {
				c.removeElem(e)
			}
		case <-timer.C:
			for e := c.queue.Back(); e != nil && time.Now().Unix() >= e.Value.(*elem).deadTime; e = e.Prev() {
				c.removeElem(e)
			}
			runtime.GC()
		case <-c.isClose:
			return
		}
	}
}
