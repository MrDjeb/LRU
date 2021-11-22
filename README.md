# LRU Cache on Go #
[Algorithm](https://ru.bmstu.wiki/LRU_(Least_Recently_Used) "Наименее недавно использованные")

***Example:***
```Go
c := NewCache(3) // {}

c.Put(1, “str1”) // {1: “str1”}

c.Put(2, “str2”) // {1: “str1”, 2: “str2”}

c.Put(3, “str3”) // {1: “str1”, 2: “str2”, 3: “str3”}

c.Get(3)

c.Get(2)

c.Get(1)

c.Get(3)

c.Put(4, “str4”) // {1: “str1”, 3: “str2”, 4: “str4”}
```

