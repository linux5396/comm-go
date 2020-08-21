package container

import (
	"fmt"
	"log"
	"testing"
)

func TestNewHashSet(t *testing.T) {
	set := NewHashSet(8)
	set1 := NewHashSet(16)
	for i := 0; i < 16; i++ {
		set.Put(i)
		if i > 10 {
			set1.Put(i)
		}
	}
	t.Log(set.Contains(12)) //true
	t.Log(set.Size())       //16
	//check
	res := set.Union(set1)
	for _, v := range res {
		fmt.Printf("%v,", v) //12,13,14,15,11,
	}
	set.Evict(12)
	t.Log(set.Contains(12)) //false
	set.Foreach(func(val interface{}) {
		log.Println("----", val)
	})
}
