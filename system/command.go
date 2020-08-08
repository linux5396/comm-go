package system

import (
	"context"
	"errors"
	"io/ioutil"
	"os/exec"
	"time"
)

//this func like C's system, can use to exec command line on linux.
//but, you should take care of this commandLine,it should not use alias, "ll" is "ls -l" ' alias.
//it impl executor is /bin/bash
func Execute(commandLine string) ([]byte, error) {
	cmd := exec.Command("/bin/bash", "-c", commandLine)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(stdout)
	return bytes, nil
}

//if run fail, will return an err,  but impl is return a txt error
//like a timeout err ,so there is a problem that how to know what's type of error
//but i think ... whatever error is,it means the command exec fail. the result is what we focus on.
func ExecuteWithTimeOut(commandLine string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", commandLine)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	//about this line , start with double-check ctx.timeout.
	//but the first checking will return a timeout err.
	//the next one is check in another go routine ,it can not return a err.
	if err := cmd.Start(); err != nil {
		return nil, err //this error is not timeout err.
	}
	//for beyond problem I post the idea below:
	if dl, _ := ctx.Deadline(); dl.Before(time.Now()) {
		return nil, errors.New("time out err")
	}
	bytes, err := ioutil.ReadAll(stdout)
	return bytes, err
}
