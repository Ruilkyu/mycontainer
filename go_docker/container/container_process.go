package container

import (
	"os"
	"os/exec"
	"syscall"
	log "github.com/sirupsen/logrus"
)

//func NewParentProcess(tty bool, command string) *exec.Cmd {
//	args := []string{"init", command}
//	cmd := exec.Command("/proc/self/exe", args...)
//	cmd.SysProcAttr = &syscall.SysProcAttr{
//		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
//	}
//	if tty{
//		cmd.Stdin = os.Stdin
//		cmd.Stdout = os.Stdout
//		cmd.Stderr = os.Stderr
//	}
//	return cmd
//}




func NewParentProcess(tty bool) (*exec.Cmd, *os.File){
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil,nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty{
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}


// 生成匿名管道
func NewPipe() (*os.File, *os.File, error){
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

