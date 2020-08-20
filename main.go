package main

import (
	//"github.com/docker/docker/pkg/reexec"
	"github.com/linux5396/comm-go/shm"
	"log"
	"os"
	"unsafe"
)

//func init() {
//	reexec.Register("cp", func() {
//
//		log.Println("start child process")
//		shmid, err := shm.Create(10289, 8, 666)
//		if err != nil {
//			log.Println(err)
//		}
//		addr, _ := shm.Attach(shmid, 0)
//		mptr := (*Data)(unsafe.Pointer(addr))
//		mptr.Code = 500
//		shm.DeAttach(addr)
//	})
//	if reexec.Init() {
//		os.Exit(0)
//	}
//}

type Data struct {
	Code int64
}

func main() {

}
