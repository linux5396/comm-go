package system

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

//register hook:customize two user signal processing methods
//can using this to impl process exits gracefully.
func RegisterSignalHandleHook(user1Sig func(), user2Sig func()) {
	go signalHandle(user1Sig, user2Sig)
}

//process control signal handler
//stop and kill signal cannot be catch
//Hook function for registering two custom reserved signals
func signalHandle(user1 func(), user2 func()) {
	defer func() {
		if r := recover(); r != nil {
			_ = fmt.Errorf("signal handle error: %v", r)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGFPE, syscall.SIGILL, syscall.SIGUSR1, syscall.SIGUSR2)
	switch <-c {
	case syscall.SIGHUP, syscall.SIGINT:
		fmt.Printf("process terminated.")
		os.Exit(0)
	case syscall.SIGQUIT:
		fmt.Printf("process quit.")
		os.Exit(0)
	case syscall.SIGFPE, syscall.SIGILL:
		fmt.Printf("process illegal exec.")
		os.Exit(-1)
	case syscall.SIGUSR1:
		user1()
	case syscall.SIGUSR2:
		user2()
	}
	signal.Stop(c)
}
