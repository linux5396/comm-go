package shm

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	//ipc flags
	IpcCreate = 01000 //IpcCreate create if key is nonexistent
	IpcExcl   = 02000 //Fail if key exists.
	IpcNoWait = 04000 //Return error on wait
	//ctl flags
	IpcRMid = 0 //Remove identifier
	IpcSet  = 1 //Set `ipc_perm' options.
	IpcStat = 2 //Get `ipc_perm' options.
	//shm attach flags
	SystemAttach   = 0x00 //系统自动选择挂载地址
	ShmReadOnly    = 0x01 //挂载只读方式,其他值为读写
	ShmRoundAttach = 0x02 //重复挂载
)

//Build shared memory based on key
func Create(key, size, mode int) (shmId uintptr, err error) {
	shmId, _, errno := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(size), uintptr(mode|IpcCreate))
	if errno != 0 {
		return 0, errors.New(errno.Error())
	}
	return shmId, nil
}

//Mount to the current process and get the address of the memory
func Attach(shmId uintptr, shmFlag int) (shmAddr uintptr, err error) {
	shmAddr, _, errno := syscall.Syscall(syscall.SYS_SHMAT, shmId, SystemAttach, uintptr(shmFlag))
	if errno != 0 {
		return 0, errors.New(errno.Error())
	}
	return shmAddr, nil
}

//unmount
func DeAttach(shmAddr uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_SHMDT, shmAddr, 0, 0)
	if errno != 0 {
		return errors.New(errno.Error())
	}
	return nil
}

//Release the shared memory to the OS
func Delete(shmId uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_SHMCTL, shmId, IpcRMid, 0)
	if errno != 0 {
		return errors.New(errno.Error())
	}
	return nil
}

//Get shared memory status
func Stat(shmId uintptr) (Shm, error) {
	shm := &Shm{}
	_, _, errno := syscall.Syscall(syscall.SYS_SHMCTL, shmId, IpcStat, uintptr(unsafe.Pointer(shm)))
	if errno != 0 {
		return *shm, errors.New(errno.Error())
	}
	return *shm, nil
}

//Underlying data structure
type Shm struct {
	IpcPerm struct {
		Key     uint32
		Uid     uint32
		Gid     uint32
		Cuid    uint32
		Cgid    uint32
		Mode    uint32
		Pad1    uint16
		Seq     uint16
		Pad2    uint16
		Unused1 uint
		Unused2 uint
	}
	ShmSegSize uint32 //段大小
	ShmATime   uint64
	ShmDTime   uint64
	ShmCTime   uint64
	ShmCPid    uint32 //最近创建该内存的进程ID
	ShmLPid    uint32 //最近使用该内存的进程ID
	ShmNAttach uint16 //挂载数
	ShmUnused  uint16
	ShmUnused2 uintptr
	ShmUnused3 uintptr
}
