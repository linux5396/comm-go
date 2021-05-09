package container

import (
	"fmt"
	"reflect"
	"sort"
)

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

type OrderMapIterator struct {
	ptr *reflect.Value
}

func NewOrderMapIterator(Map interface{}) *OrderMapIterator {
	v := reflect.ValueOf(Map)
	if v.Kind() != reflect.Map {
		panic("data is not map")
	}
	return &OrderMapIterator{&v}
}

//实现map的指定顺序迭代，如果返回true，迭代是符合要求的；返回false则未定义错误
func (it *OrderMapIterator) Foreach(do func(key, val interface{}), less func(k1, k2 interface{}) bool) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("iterator failed:", r)
		}
	}()
	keys := make([]reflect.Value, 0)
	for _, key := range it.ptr.MapKeys() {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return less(keys[i].Interface(), keys[j].Interface())
	})
	for i := 0; i < len(keys); i++ {
		do(keys[i].Interface(), it.ptr.MapIndex(keys[i]).Interface())
	}
	return true
}

//以下实现了一些常用的比较器
func IntLess(k1, k2 interface{}) bool {
	n1 := k1.(int)
	n2 := k2.(int)
	if n1 < n2 {
		return true
	}
	return false
}

func Float32Less(k1, k2 interface{}) bool {
	n1 := k1.(float32)
	n2 := k2.(float32)
	if n1 < n2 {
		return true
	}
	return false
}

func Float64Less(k1, k2 interface{}) bool {
	n1 := k1.(float64)
	n2 := k2.(float64)
	if n1 < n2 {
		return true
	}
	return false
}

func Int32Less(k1, k2 interface{}) bool {
	n1 := k1.(int32)
	n2 := k2.(int32)
	if n1 < n2 {
		return true
	}
	return false
}

func Int64Less(k1, k2 interface{}) bool {
	n1 := k1.(int64)
	n2 := k2.(int64)
	if n1 < n2 {
		return true
	}
	return false
}
