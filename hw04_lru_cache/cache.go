package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Clear() {
	for {
		if element := c.queue.Back(); element != nil {
			c.queue.Remove(element)
			item := element.Value.(*cacheItem)
			delete(c.items, item.key)
		}
		if c.queue.Front() == nil {
			break
		}
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if element, exists := c.items[key]; exists {
		c.queue.MoveToFront(element)
		element.Value.(*cacheItem).value = value
		return true
	}

	if c.queue.Len() == c.capacity {
		if element := c.queue.Back(); element != nil {
			c.queue.Remove(element)
			item := element.Value.(*cacheItem)
			delete(c.items, item.key)
		}
	}

	item := &cacheItem{
		key:   key,
		value: value,
	}

	element := c.queue.PushFront(item)
	c.items[item.key] = element

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	element, exists := c.items[key]
	if !exists {
		return nil, false
	}
	c.queue.MoveToFront(element)
	return element.Value.(*cacheItem).value, true
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
