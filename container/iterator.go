package container

import "reflect"

//iterator interface
type Foreach interface {
	Foreach(f func(arg interface{})) Foreach
}

//slice iterator
//support function style
type SliceIterator struct {
	ptr *reflect.Value
}

//create a new slice iterator
//the parameter must be a slice otherwise panic
func NewSliceIterator(data interface{}) *SliceIterator {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		panic("data is not slice")
	}
	return &SliceIterator{ptr: &v}
}

//iterative interface
func (it *SliceIterator) Foreach(f func(itf interface{})) Foreach {
	if it.ptr.Len() > 0 {
		if !it.ptr.Index(0).CanInterface() {
			panic("this slice does not support interface")
		}
	}
	for i := 0; i < it.ptr.Len(); i++ {
		f(it.ptr.Index(i).Interface())
	}
	return it
}
