package lru_cache

import (
	"container/list"
	"errors"
)

//non Synchronize dequeue
//impl by list.List
type Dequeue struct {
	l *list.List
}

var emptyErr = errors.New("queue is empty")

func NewDequeue() *Dequeue {
	l := list.New()
	l.Init()
	return &Dequeue{l}
}

func (d *Dequeue) PushBack(val interface{}) {
	d.l.PushBack(val)
}

func (d *Dequeue) PushFirst(val interface{}) {
	d.l.PushFront(val)
}

func (d *Dequeue) PollLast() (interface{}, error) {
	if d.l.Len() == 0 {
		return nil, emptyErr
	}
	e := d.l.Back()
	d.l.Remove(e)
	return e.Value, nil
}

func (d *Dequeue) PollFirst() (interface{}, error) {
	if d.l.Len() == 0 {
		return nil, emptyErr
	}
	e := d.l.Front()
	d.l.Remove(e)
	return e.Value, nil
}

func (d *Dequeue) PeekLast() (interface{}, error) {
	if d.l.Len() == 0 {
		return nil, emptyErr
	}
	e := d.l.Back()
	return e.Value, nil
}

func (d *Dequeue) PeekFirst() (interface{}, error) {
	if d.l.Len() == 0 {
		return nil, emptyErr
	}
	e := d.l.Front()
	return e.Value, nil
}

func (d *Dequeue) Clear() {
	d.l.Init()
}

func (d *Dequeue) Size() int {
	return d.l.Len()
}
