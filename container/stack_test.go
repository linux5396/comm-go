package container

import "testing"

type Holder struct {
	Stack
}

func TestNonSynchronizeStack(t *testing.T) {
	h := Holder{GetNonSynchronizeStack(8)}
	h.Push(1)
	t.Log(h.Size())  //1
	t.Log(h.Peek())  //1
	t.Log(h.Empty()) //false
	t.Log(h.Pop())
	t.Log(h.Full())
	for i := 0; i < 9; i++ {
		t.Log(h.Push(i))
	}
}
