package main

import (
	"fmt"
	"sync"
)

func main() {
	c := newCache(2)
	c.put("a", 1)
	c.put("a", 2)
	c.put("b", 3)
	c.put("c", 4)

	fmt.Println(c.get("b"))
	fmt.Printf("%+v", c)

}

type doublyLinkedList struct {
	head   *node
	tail   *node
	length int
}

type node struct {
	next  *node
	prev  *node
	value value
}

type value struct {
	key   string
	value int
}

func (l *doublyLinkedList) push(e *node) {
	if e == nil {
		return
	}

	l.length++

	if l.head == nil {
		l.head = e
		l.tail = e
		return
	}

	l.tail.next = e
	e.prev = l.tail
	l.tail = e
}

func (l *doublyLinkedList) delete(e *node) {
	if e == nil || l.length == 0 {
		return
	}

	l.length--

	if l.head == e {
		l.head = e.next
	}

	if l.tail == e {
		l.tail = e.prev
	}

	if e.next != nil {
		e.next = e.next.next
	}

	if e.prev != nil {
		e.prev = e.prev.prev
	}

	e.next = nil
	e.prev = nil
}

type cache struct {
	cache    map[string]*node
	queue    doublyLinkedList
	capacity int
	mu       sync.RWMutex
}

func newCache(capacity int) *cache {
	c := make(map[string]*node, capacity)
	return &cache{cache: c, capacity: capacity}
}

func (c *cache) put(key string, newValue int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	storedNode, ok := c.cache[key]
	if ok {
		c.queue.delete(storedNode)
	}

	newNode := node{value: value{key: key, value: newValue}}
	c.cache[key] = &newNode
	c.queue.push(&newNode)

	if c.queue.length > c.capacity {
		firstNode := c.queue.head
		c.queue.delete(firstNode)
		delete(c.cache, firstNode.value.key)
	}
}

func (c *cache) get(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	storedNode, ok := c.cache[key]
	if !ok {
		return 0
	}

	c.queue.delete(storedNode)
	c.queue.push(storedNode)

	return storedNode.value.value
}

func (c *cache) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	storedNode, ok := c.cache[key]
	if !ok {
		return
	}

	c.queue.delete(storedNode)
	delete(c.cache, key)
}
