package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mx       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	item, exists := c.items[key]
	if exists {
		datum := item.Value.(cacheItem)
		datum.value = value
		item.Value = datum

		c.queue.MoveToFront(item)
	} else {
		item = c.queue.PushFront(cacheItem{key: key, value: value})
	}

	if c.queue.Len() > c.capacity {
		if back := c.queue.Back(); back != nil {
			c.queue.Remove(back)

			delete(c.items, back.Value.(cacheItem).key)
		}
	}

	if c.capacity > 0 {
		c.items[key] = item
	}

	return exists
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if val, ok := c.items[key]; ok {
		result := val.Value.(cacheItem)
		c.queue.MoveToFront(val)

		return result.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
