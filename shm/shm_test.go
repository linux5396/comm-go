package shm

import (
	"github.com/linux5396/comm-go/system"
	"testing"
	"time"
	"unsafe"
)

type DataPack struct {
	code   int
	number float64
}

func TestNew(t *testing.T) {
	h := New(10086)
	data := DataPack{}
	addr, err := h.GetShm(data, false)
	if err != nil {
		t.Error(err)
	}
	ptr := (*DataPack)(unsafe.Pointer(addr))
	ptr.code = 200
	t.Logf("%+v", *ptr)
	_ = h.DestroyShm()
	_, _ = system.Execute(`ls -a`)
	bf, err1 := system.ExecuteWithTimeOut(`sleep 10`, time.Millisecond*200000)
	if err1 != nil {
		t.Error(err1)
	} else {
		t.Log(string(bf))
	}
}
