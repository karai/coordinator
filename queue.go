package main

import (
	"fmt"
	"sync"
)

type workQueue struct {
	stack []string
	lock  sync.RWMutex
}

func (c *workQueue) push(object string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.stack = append(c.stack, object)
}

func (c *workQueue) pop() error {
	len := len(c.stack)
	if len > 0 {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.stack = c.stack[:len-1]
		return nil
	}
	return fmt.Errorf("Pop Error: Queue is empty")
}

func (c *workQueue) front() (string, error) {
	len := len(c.stack)
	if len > 0 {
		c.lock.Lock()
		defer c.lock.Unlock()
		return c.stack[len-1], nil
	}
	return "", fmt.Errorf("Peep Error: Queue is empty")
}

func (c *workQueue) size() int {
	return len(c.stack)
}

func (c *workQueue) empty() bool {
	return len(c.stack) == 0
}

func queue() {
	customQueue := &workQueue{
		stack: make([]string, 0),
	}
	fmt.Printf("Push: A\n")
	customQueue.push("A")
	fmt.Printf("Pop: B\n")
	customQueue.push("B")
	fmt.Printf("Size: %d\n", customQueue.size())
	for customQueue.size() > 0 {
		frontVal, _ := customQueue.front()
		fmt.Printf("Front: %s\n", frontVal)
		fmt.Printf("Pop: %s\n", frontVal)
		customQueue.pop()
	}
	fmt.Printf("Size: %d\n", customQueue.size())
}
