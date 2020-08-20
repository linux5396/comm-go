package container

import (
	"errors"
	"sync"
)

//define a group method of stack operations.
//type interface to make it abstract to Support better scalability
type Stack interface {
	Push(val interface{}) error
	Peek() (interface{}, error)
	Pop() (interface{}, error)
	Clear()
	Size() int
	Empty() bool
	Full() bool
}

//non Synchronize stack
type NonSynchronizeStack struct {
	len   int
	max   int
	slice []interface{}
}

func GetNonSynchronizeStack(max int) *NonSynchronizeStack {
	if max < 1 {
		panic("max can not lt 1")
	}
	return &NonSynchronizeStack{
		max: max,
	}
}

func (n *NonSynchronizeStack) Push(val interface{}) error {
	if n.len == n.max {
		return errors.New("stack is overflow")
	}
	n.slice = append(n.slice, val)
	n.len++
	return nil
}

func (n *NonSynchronizeStack) Peek() (interface{}, error) {
	if n.len > 0 {
		return n.slice[n.len-1], nil
	} else {
		return nil, errors.New("stack is empty")
	}
}

func (n *NonSynchronizeStack) Pop() (interface{}, error) {
	if n.len > 0 {
		val := n.slice[n.len-1]
		n.len--
		n.slice = n.slice[0:n.len]
		return val, nil
	} else {
		return nil, errors.New("stack is empty")
	}
}

func (n *NonSynchronizeStack) Clear() {
	n.len = 0
	n.slice = make([]interface{}, 0)
}

func (n *NonSynchronizeStack) Size() int {
	return n.len
}

func (n *NonSynchronizeStack) Empty() bool {
	return n.len == 0
}
func (n *NonSynchronizeStack) Full() bool {
	return n.len == n.max
}

//base on normal stack, but with rw-lock to keep safe
type ConcurrentStack struct {
	sync.RWMutex //Use read-write locks to implement the stack
	inner        NonSynchronizeStack
}

func NewConcurrentStack(max int) *ConcurrentStack {
	return &ConcurrentStack{
		RWMutex: sync.RWMutex{},
		inner:   *GetNonSynchronizeStack(max),
	}
}

func (c *ConcurrentStack) Push(val interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.inner.Push(val)
}

func (c *ConcurrentStack) Peek() (interface{}, error) {
	c.RLock()
	defer c.RUnlock()
	return c.inner.Peek()
}

func (c *ConcurrentStack) Pop() (interface{}, error) {
	c.Lock()
	defer c.Unlock()
	return c.inner.Pop()
}

func (c *ConcurrentStack) Clear() {
	c.Lock()
	defer c.Unlock()
	c.inner.Clear()
}

func (c *ConcurrentStack) Size() int {
	c.Lock()
	defer c.Unlock()
	return c.inner.Size()
}

func (c *ConcurrentStack) Empty() bool {
	c.Lock()
	defer c.Unlock()
	return c.inner.Empty()
}

func (c *ConcurrentStack) Full() bool {
	c.Lock()
	defer c.Unlock()
	return c.inner.Full()
}
