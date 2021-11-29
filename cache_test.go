package cache

import (
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"
)

func TestCache1(t *testing.T) {
	runtime.GOMAXPROCS(4)
	c, err := NewCache(5)
	if err != nil {
		t.Fatal("err should be nil")
	}
	defer c.Destroy()

	fmt.Printf("TTL: %v, checkTTL: %v, curSise: %v, minSize: %v, maxSize: %v",
		c.ttl, c.checkTTL, c.curSize, c.minSize, c.maxSize)

	var someKey KeyT = 0
	c.Put(someKey, "test")
	testVal, testOk := c.Get(someKey)
	if testVal == "" {
		t.Fatal("value should not be nil")
	}
	if !testOk {
		t.Fatal("ok shoud be true")
	}

	time.Sleep(3 * time.Second)
	testVal, testOk = c.Get(someKey)
	if testVal == "" {
		t.Fatal("value should not be nil, timer is bad")
	}
	if !testOk {
		t.Fatal("ok shoud be true, timer is bad")
	}

	testVal, testOk = c.Get(someKey + 1)
	if testVal != "" {
		t.Fatal("value should be nil")
	}
	if testOk {
		t.Fatal("ok shoud be false")
	}

	var testUpdate ValueT = "is"
	c.Put(someKey, testUpdate)
	testVal, testOk = c.Get(someKey)
	if testVal != testUpdate {
		t.Fatal("update value is bad")
	}
	c.removeElem(c.queue.Front())
	if _, ok := c.Get(someKey); ok {
		t.Fatal("removeElem() is bad")
	}

	c.Display()

	time.Sleep(3 * time.Second)
	if c.curSize != 0 {
		t.Fatalf("c.curSize() should be 0, but it %v, calculate c.curSize in <-timer.C is bad", c.curSize)
	}

	c.Put(1, "test1")
	c.Put(2, "test2")
	c.Put(3, "test3")

	for e := c.queue.Back(); e != nil; e = e.Prev() {
		fmt.Println("First = ", e.Value.(*elem))
	}
	c.Get(2)

	for e := c.queue.Back(); e != nil; e = e.Prev() {
		fmt.Println("Second = ", e.Value.(*elem))
	}
}

func BenchmarkPut(b *testing.B) {
	c, err := NewCache(5)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Destroy()

	keys := make([]KeyT, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = KeyT((i + 4567) % 2345)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Put(keys[i], "test")
	}
}
