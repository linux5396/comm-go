package lru_cache

import "testing"

func TestNewDequeue(t *testing.T) {
	dequeue := NewDequeue()
	t.Log(dequeue.Size()) //0
	dequeue.PushBack("hello")
	dequeue.PushBack("hello1")
	dequeue.PushBack("hello2")
	t.Log(dequeue.Size())      //3
	t.Log(dequeue.PollFirst()) //hello
	t.Log(dequeue.PollFirst()) //hello1
	t.Log(dequeue.PollFirst()) //hello2
}
