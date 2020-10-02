package system

import (
	"context"
	"errors"
	"io/ioutil"
	"os/exec"
	"time"
)

//using customize parser to execute cmd line.
func ExecuteWithParser(cmdParse, commandLine string) (out []byte, outErr []byte, err error) {
	cmd := exec.Command(cmdParse, "-c", commandLine)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, nil, err
	}
	//read stdout and stderr stream from
	outErr, err = ioutil.ReadAll(stderr)
	out, err = ioutil.ReadAll(stdout)
	return
}

//this func like C's system, can use to exec command line on linux.
//but, you should take care of this commandLine,it should not use alias, "ll" is "ls -l" ' alias.
//it impl executor is /bin/bash
func Execute(commandLine string) (out []byte, outErr []byte, err error) {
	cmd := exec.Command("/bin/bash", "-c", commandLine)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, nil, err
	}
	//read stdout and stderr stream from
	outErr, err = ioutil.ReadAll(stderr)
	out, err = ioutil.ReadAll(stdout)
	return
}

//if run fail, will return an err,  but impl is return a txt error
//like a timeout err ,so there is a problem that how to know what's type of error
//but i think ... whatever error is,it means the command exec fail. the result is what we focus on.
//this func only support the
func ExecuteWithTimeOut(commandLine string, timeout time.Duration) (out []byte, outErr []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", commandLine)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	//about this line , start with double-check ctx.timeout.
	//but the first checking will return a timeout err.
	//the next one is check in another go routine ,it can not return a err.
	if err := cmd.Start(); err != nil {
		return nil, nil, err //this error is not timeout err.
	}
	//for beyond problem I post the idea below:
	//but ..because the timeout ctx without a way to get the timeout err,only by the blocked chan.
	//and our application should run in a non-blocked way,so we only call the deadline  manually.
	//this way exists a error:  There are nano scale errors.
	//thus if you want without a mistake to judge time out or not.
	//you should enhance the exec commandline,like get a file's locker or get a callback result.
	if dl, _ := ctx.Deadline(); dl.Before(time.Now()) {
		return nil, nil, errors.New("time out err")
	}
	out, err = ioutil.ReadAll(stdout)
	outErr, err = ioutil.ReadAll(stderr)
	return
}
