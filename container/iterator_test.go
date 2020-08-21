package container

import (
	"fmt"
	"testing"
)

type Man struct {
	val int
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
}
