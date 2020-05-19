package container

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
	_ "time"
)


var (
	RUNNING                string = "running"
	STOP                   string = "stopped"
	Exit                   string = "exited"
	DefaultInfoLocation    string = "/var/run/mycontainer/%s/"
	ConfigName             string = "config.json"
	ContainerLogFile	   string = "container.log"
	RootUrl                string = "/root"
	MntUrl                 string = "/root/mnt/%s"
	WriteLayerUrl          string = "/root/writeLayer/%s"
	ReadLayerUrl           string = "/root/readLayer/%s"
	WorkLayerUrl           string = "/root/workLayer/%s"
)

type ContainerInfo struct{
	Pid           string `json:"pid"`  //容器init进程在宿主机的PID
	Id            string `json:"id"`
	Name          string `json:"name"`
	Command       string `json:"command"` //容器init的运行命令
    CreatedTime   string `json:"createdTime"`
	Status        string `json:"status"`
	Volume        string `json:"volume"`
}

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




func NewParentProcess(tty bool, volume string, containerName string, logfile bool, imageName string, envSlice []string) (*exec.Cmd, *os.File){
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil,nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirURL, 0622); err != nil {
		log.Errorf("NewParentProcess mkdir %s error %v", dirURL, err)
		return nil, nil
	}
	stdLogFilePath := dirURL + ContainerLogFile
	stdLogFile, err := os.Create(stdLogFilePath)
	if err != nil {
		log.Errorf("NewParentProcess create file %s error %v", stdLogFilePath, err)
		return nil, nil
	}

	if tty  && !logfile{
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else if tty  && logfile{
		cmd.Stdin = os.Stdin
		cmd.Stdout = stdLogFile
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Env = append(os.Environ(), envSlice...)
	NewWorkSpace(volume, imageName, containerName)
	cmd.Dir = fmt.Sprintf(MntUrl, containerName)
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






