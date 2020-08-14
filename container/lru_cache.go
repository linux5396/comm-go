package lru_cache

import (
	"container/list"
	"fmt"
	"sync"
)

//Support concurrent and safe caching
//this version uses map and mutex to implement.
//In the future，I will make bench test for a sync.Map to choose which performance is better.
type LRUCache struct {
	sync.Mutex
	capacity int                           // maximum number of key-value pairs
	cache    map[interface{}]*list.Element // map for cached key-value pairs
	lru      *list.List                    // LRU list
}

type Pair struct {
	key   interface{} // cache key
	value interface{} // cache value
}

// NewLRUCache returns a new, empty LRUCache
func NewLRUCache(capacity int) *LRUCache {
	c := new(LRUCache)
	c.capacity = capacity
	c.cache = make(map[interface{}]*list.Element)
	c.lru = list.New()
	return c
}

// Get get cached value from LRU cache
// The second return value indicates whether key is found or not, true if found, false if not
func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.lru.MoveToFront(elem) // move node to head of lru list
		return elem.Value.(*Pair).value, true
	}
	return nil, false
}

// Add adds a key-value pair to LRU cache, true if eviction occurs, false if not
func (c *LRUCache) Add(key interface{}, value interface{}) bool {
	c.Lock()
	defer c.Unlock()
	// update item if found in cache
	if elem, ok := c.cache[key]; ok {
		c.lru.MoveToFront(elem) // update lru list
		elem.Value.(*Pair).value = value
		return false
	}
	// add item if not found
	elem := c.lru.PushFront(&Pair{key, value})
	c.cache[key] = elem
	// evict item if needed
	if c.lru.Len() > c.capacity {
		c.evict()
		return true
	}
	return false
}

// evict a key-value pair from LRU cache
func (c *LRUCache) evict() {
	elem := c.lru.Back()
	if elem == nil {
		return
	}
	// remove item at the end of lru list
	c.lru.Remove(elem)
	delete(c.cache, elem.Value.(*Pair).key)
}

// Del deletes cached value from cache
func (c *LRUCache) Del(key interface{}) {
	c.Lock()
	defer c.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.lru.Remove(elem)
		delete(c.cache, key)
	}
}

// Len returns number of items in cache
func (c *LRUCache) Len() int {
	c.Lock()
	defer c.Unlock()
	return c.lru.Len()
}

// Keys returns keys of items in cache
func (c *LRUCache) Keys() []interface{} {
	var keyList []interface{}
	c.Lock()
	for key := range c.cache {
		keyList = append(keyList, key)
	}
	c.Unlock()
	return keyList
}

//EnlargeCapacity enlarges the capacity of cache
func (c *LRUCache) EnlargeCapacity(newCapacity int) error {
	// lock
	c.Lock()
	defer c.Unlock()
	// check newCapacity
	if newCapacity < c.capacity {
		return fmt.Errorf("newCapacity[%d] must be larger than current[%d]",
			newCapacity, c.capacity)
	}
	c.capacity = newCapacity
	return nil
}
