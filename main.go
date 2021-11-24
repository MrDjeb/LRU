package main

/*
	Implement:
  		* memory cache size,
  		* TLL,
  		* testing,
  		* sync.mutex/sync.semaphore
*/

import (
	"github.com/MrDjeb/LRU/cache"
)

func main() {
	c, _ := cache.NewCache(3)
	c.Put(1, "str1")
	c.Put(2, "str2")
	c.Put(3, "str3")

	c.Get(3)
	c.Get(2)
	c.Get(1)
	c.Get(3)
	c.Put(4, "str4")
	c.Display()
}
