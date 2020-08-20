package system

import (
	"testing"
	"time"
)

func TestPoppyLock_LockWithTimeout(t *testing.T) {
	poppy := PoppyLock{}
	go func() {
		poppy.Lock()
		time.Sleep(time.Second * 3)
		poppy.Unlock()
	}()
	time.Sleep(time.Millisecond)
	go func() {
		ok := poppy.TryLock(time.Second * 1)
		t.Log(ok)
		if ok {
			poppy.Unlock()
		}
	}()
	time.Sleep(time.Second * 5)
}
