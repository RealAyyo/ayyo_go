package hw04lrucache

import "sync"

type Key string

type ListValue struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.items[key]; ok {
		node.Value = ListValue{key: key, value: value}
		c.queue.MoveToFront(node)
		return true
	}

	if c.queue.Len() == c.capacity {
		removedKey := c.queue.Back().Value.(ListValue).key
		c.queue.Remove(c.queue.Back())
		delete(c.items, removedKey)
	}

	c.items[key] = c.queue.PushFront(ListValue{key: key, value: value})

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.items[key]; ok {
		c.queue.MoveToFront(node)
		return node.Value.(ListValue).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[Key]*ListItem)
	c.queue = NewList()
}
