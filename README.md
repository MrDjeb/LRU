# LRU Cache on Go #
[Algorithm](https://ru.bmstu.wiki/LRU_(Least_Recently_Used) "Наименее недавно использованные")

***Example:***
```Go
package main

import (
	"log"
	"time"

	"github.com/MrDjeb/LRU/cache"
)

func main() {
	c, err := cache.NewCache(5)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Destroy()
	c.Put(1, "str1")
	c.Put(2, "str2")
	c.Put(3, "str3")
	c.Get(3)
	c.Get(2)
	c.Get(1)
	c.Get(3)
	c.Put(4, "str4")
	c.Display()
	time.Sleep(6 * time.Second)
	c.Display()
}
```

