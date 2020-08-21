Comm golang library

## shm

> # 	shm is a comm share memory lib.
>
> > ​	Now It supports OS on Linux. You can use like below:
> >
> > ```go
> > type DataPack struct {
> >    code   int
> >    number float64
> > }
> > 
> > func TestNew(t *testing.T) {
> >    h := New(10086)
> >    data := DataPack{}
> >    addr, err := h.GetShm(data, false)
> >    if err != nil {
> >       t.Error(err)
> >    }
> >    ptr := (*DataPack)(unsafe.Pointer(addr))
> >    ptr.code = 200
> >    t.Logf("%+v", *ptr)
> >    _ = h.DestroyShm()
> > }
> > ```
>
> > yeah, It wrap so many details. you can new a holder, and get a memory segment by specifying the shm key. And you can pass the arg of readOnly or not to determine this shm's feature.
> >
> > shm is using on ipc, like async log、async task ,etc.
>
> ----

## container

> Container package contains some public containers,like queues, collections, LRU caches, stacks, etc., while providing slice iterator.
> some examples of use are provided below：
> >
> > ```go
> >     //SET
> > 	set := NewHashSet(8)
> > 	set1 := NewHashSet(16)
> > 	for i := 0; i < 16; i++ {
> > 		set.Put(i)
> > 		if i > 10 {
> > 			set1.Put(i)
> > 		}
> > 	}
> > 	t.Log(set.Contains(12)) //true
> > 	t.Log(set.Size())       //16
> > 	//check
> > 	res := set.Union(set1)
> > 	for _, v := range res {
> > 		fmt.Printf("%v,", v) //12,13,14,15,11,
> > 	}
> > 	set.Evict(12)
> > 	t.Log(set.Contains(12)) //false
> > 	set.Foreach(func(val interface{}) {
> > 		log.Println("----", val)
> > 	})
> > ```
> >
> > 

