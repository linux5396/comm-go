package shm

import (
	"testing"
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
}
