# Comm golang library

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
> 