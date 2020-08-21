package system

import (
	"sync"
	"testing"
	"time"
)

func TestPoppyLock_LockWithTimeout(t *testing.T) {
	poppy := PoppyLock{}
	//go func() {
	//	poppy.Lock()
	//	time.Sleep(time.Second * 3)
	//	poppy.Unlock()
	//}()
	//time.Sleep(time.Millisecond)
	qg := sync.WaitGroup{}
	qg.Add(1)
	for i := 0; i < 10; i++ {
		go func() {
			qg.Wait()
			ok := poppy.TryLock(time.Millisecond * 10)
			t.Log(ok)
			if ok {
				time.Sleep(time.Millisecond * 8)
				poppy.Unlock()
			}
		}()
	}
	qg.Done()
	time.Sleep(time.Second * 5)
}
