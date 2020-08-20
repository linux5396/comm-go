package shm

import "unsafe"

//Shared memory holder
type ShareMemoryHolder struct {
	shmId   uintptr
	shmAddr uintptr
	ShmKey  int
}

//a key range from 0 to 65535
func New(shmKey uint16) *ShareMemoryHolder {
	return &ShareMemoryHolder{
		ShmKey: int(shmKey),
	}
}

//new by transform func and a key of T.
//f should return a uint which range is from 0 to 65535
//so this func can be a hash func or
func NewByFunc(key interface{}, f func(key interface{}) uint16) *ShareMemoryHolder {
	return &ShareMemoryHolder{
		ShmKey: int(f(key)),
	}
}

type ShareMemoryItf interface {
	//根据Key获取并挂载到当前进程,基于共享结构来申请内存，可将申请到的指针转为结构，来进行操作
	GetShm(shareStruct interface{}, readOnly bool) (uintptr, error)
	//直接申请一块大小,暂不支持IO接口，因此使用这个，仍然建议转为结构体指针进行操作
	GetShmInSize(size int, readOnly bool) (uintptr, error)
	//解除挂载
	DeAttach() error
	//彻底删除，只由一个对象执行
	DestroyShm() error
}

func (s *ShareMemoryHolder) GetShm(shareStruct interface{}, readOnly bool) (uintptr, error) {
	shmid, err := Create(s.ShmKey, int(unsafe.Sizeof(shareStruct)), 0666)
	if err != nil {
		return 0, err
	}
	if readOnly {
		s.shmAddr, err = Attach(shmid, ShmReadOnly)
	} else {
		s.shmAddr, err = Attach(shmid, 0)
	}
	if err != nil {
		return 0, err
	}
	s.shmId = shmid
	return s.shmAddr, nil
}

func (s *ShareMemoryHolder) GetShmInSize(size int, readOnly bool) (uintptr, error) {
	if size < 0 {
		panic("size can not lt 0")
	}
	shmid, err := Create(s.ShmKey, size, 0666)
	if err != nil {
		return 0, err
	}
	if readOnly {
		s.shmAddr, err = Attach(shmid, ShmReadOnly)
	} else {
		s.shmAddr, err = Attach(shmid, 0)
	}
	if err != nil {
		return 0, err
	}
	s.shmId = shmid
	return s.shmAddr, nil
}

func (s *ShareMemoryHolder) DeAttach() error {
	return DeAttach(s.shmAddr)
}

func (s *ShareMemoryHolder) DestroyShm() error {
	_, err := Stat(s.shmId)
	if err != nil {
		return err
	}
	return Delete(s.shmId)
}

//TODO support IO interface
