package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)


//func RunContainerInitProcess(command string, args []string) error {
//	logrus.Infof("command %s", command)
//
//	syscall.Mount("", "/", "", syscall.MS_PRIVATE | syscall.MS_REC, "")
//
//	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
//	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
//	argv := []string{command}
//	if err:= syscall.Exec(command, argv, os.Environ()); err != nil {
//		logrus.Errorf(err.Error())
//	}
//	return nil
//}


func readUserCommand() []string{
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil{
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

func RunContainerInitProcess() error{
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0{
		return fmt.Errorf("run container get user command error, cmdArray is nil")
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	path, err := exec.LookPath(cmdArray[0])
	if err != nil{
		log.Errorf("exec loop path error %v", err)
		return err
	}
	log.Infof("find path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil{
		log.Errorf(err.Error())
	}
	return nil
}
