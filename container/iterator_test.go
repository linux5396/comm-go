package container

import (
	"fmt"
	"testing"
	"unsafe"
)

type Man struct {
	val int
}

func memset(slicePtr unsafe.Pointer, c byte, n uintptr) {
	ptr := uintptr(slicePtr)
	var i uintptr
	for i = 0; i < n; i++ {
		pByte := (*byte)(unsafe.Pointer(ptr + i))
		*pByte = c
	}
}

func TestSlice(t *testing.T) {
	NewSliceIterator([]Man{{val: 1}, {val: 2}, {val: 3}, {val: 4}}).Foreach(func(arg interface{}) {
		t.Log(arg.(Man).val)
	}).Foreach(func(itf interface{}) {
		if itf.(Man).val > 2 {
			fmt.Println(itf.(Man).val)
		}
	}).Foreach(func(arg interface{}) {
		//something call
	}).Foreach(func(arg interface{}) {
		//something call
	}).Foreach(func(arg interface{}) {
		//something call
	})

	ff := make([]float64, 12)
	memset(unsafe.Pointer(&ff), 1.0, unsafe.Sizeof(ff))
	t.Log(ff[0])

}

func TestOrderMapIterator_Foreach(t *testing.T) {
	m := map[int]string{
		1:  "a",
		2:  "b",
		3:  "c",
		-1: "dd",
	}
	it := NewOrderMapIterator(m)
	it.Foreach(func(key, val interface{}) {
		if k, ok := key.(int); ok {
			if v, ok2 := val.(string); ok2 {
				fmt.Println(k, ":", v)
			}
		}
	}, func(k1, k2 interface{}) bool {
		return IntLess(k2, k1) //反向
	})
}
