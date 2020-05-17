package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	log "github.com/Sirupsen/logrus"
	_"time"
)


var (
	RUNNING                string = "running"
	STOP                   string = "stopped"
	Exit                   string = "exited"
	DefaultInfoLocation    string = "/var/run/mycontainer/%s/"
	ConfigName             string = "config.json"
	ContainerLogFile	   string = "container.log"
)

type ContainerInfo struct{
	Pid           string `json:"pid"`  //容器init进程在宿主机的PID
	Id            string `json:"id"`
	Name          string `json:"name"`
	Command       string `json:"command"` //容器init的运行命令
    CreatedTime   string `json:"createdTime"`
	Status        string `json:"status"`
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




func NewParentProcess(tty bool, volume string, containerName string, logfile bool) (*exec.Cmd, *os.File){
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
	mntURL := "/root/merged/"
	rootURL := "/root/"
	NewWorkSpace(rootURL, mntURL, volume)
	cmd.Dir = mntURL
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

// 创建一个Overlay系统作为容器的根目录
func NewWorkSpace(rootURL string, mntURL string, volume string){
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateWorkLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)



	// 判断是否执行挂载数据卷操作
	if volume != ""{
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(rootURL, mntURL, volumeURLs)
			log.Infof("%q", volumeURLs)
		} else {
			log.Infof("Volume parameter input is not correct.")
		}
	}
}


// 解析volume字符串
func volumeUrlExtract(volume string) ([]string) {
	var volumeURLs []string
	volumeURLs = strings.Split(volume, ":")
	return volumeURLs
}


// 挂载数据卷
func MountVolume(rootURL string, mntURL string, volumeURLS []string) {
	// 在宿主机常见宿主目录
	parentUrl := volumeURLS[0]
	if err := os.Mkdir(parentUrl, 0777); err != nil {
		log.Infof("Mkdir parent dir %s error. %v", parentUrl, err)
	}
	CreateVolumePoint(parentUrl)

	// 在容器里创建挂载目录
	containerUrl := volumeURLS[1]
	containerVolumeURL := mntURL + containerUrl
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
		log.Infof("Mkdir container dir %s error. %v", containerVolumeURL, err)
	}

	// 把宿主机文件目录挂载到容器挂载点
	lower := "lowerdir=" + parentUrl + "/readLayer/"
	upper := "upperdir=" + parentUrl + "/writeLayer/"
	work := "workdir=" + parentUrl + "/workLayer/"
	parm := lower + "," + upper + "," + work
	//cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", parm, containerVolumeURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Mount volume failed. %v", err)
	}
}



func CreateVolumePoint(parentUrl string) {
	readURL := parentUrl + "/readLayer/"
	writeURL := parentUrl + "/writeLayer/"
	workURL := parentUrl + "/workLayer/"

	if err := os.Mkdir(readURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", readURL, err)
	}

	if err := os.Mkdir(writeURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeURL, err)
	}

	if err := os.Mkdir(workURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", workURL, err)
	}
}



func CreateReadOnlyLayer(rootURL string){
	busyboxURL := rootURL + "busybox/"
	busyboxTarURL := rootURL + "busybox.tar"
	exist, err := PathExists(busyboxURL)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", busyboxURL, err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error. %v", busyboxURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", busyboxURL, err)
		}
	}
}

func CreateWorkLayer(rootURL string) {
	workURL := rootURL + "workLayer/"
	if err := os.Mkdir(workURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", workURL, err)
	}
}


func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeURL, err)
	}
}


func CreateMountPoint(rootURL string, mntURL string) {
	if err := os.Mkdir(mntURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mntURL, err)
	}
	lower := "lowerdir=" + rootURL + "busybox"
	upper := "upperdir=" + rootURL + "writeLayer"
	work := "workdir=" + rootURL + "workLayer"
	parm := lower + "," + upper + "," + work
	//cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", parm, mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}

// 当容器退出时删除Overlay系统
func DeleteWorkSpace(rootURL string, mntURL string, volume string){
	if (volume != ""){
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if (length == 2 && volumeURLs[0] != "" && volumeURLs[1] != ""){
			DeleteMountPointWithVolume(rootURL, mntURL, volumeURLs)
		} else {
			DeleteMountPoint(rootURL, mntURL)
		}
	} else {
		DeleteMountPoint(rootURL, mntURL)
	}
	DeleteWriteLayer(rootURL)
	DeleteWorkLayer(rootURL)
}


func DeleteMountPointWithVolume(rootURL string, mntURL string, volumeURLs []string){

	// 卸载容器里的volume挂载点
	containerUrl := mntURL + volumeURLs[1]
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Umount volume failed. %v", err)
	}

	// 卸载整个容器文件系统挂载点
	cmd = exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Umount mountpoint failed. %v", err)
	}else{
		log.Infof("Umount mountpoint success")
	}
	// 删除文件系统挂载点
	if err := os.RemoveAll(mntURL); err != nil {
		log.Infof("Remove mountpoint dir %s error %v", mntURL, err)
	}
}



func DeleteMountPoint(rootURL string, mntURL string){
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout=os.Stdout
	cmd.Stderr=os.Stderr
	a := exec.Command("umount", mntURL)
	a.Run()
	if err := cmd.Run(); err != nil {
		log.Errorf("%v",err)
	}else{
		log.Infof("Umount mountpoint success")
	}

	//time.Sleep(time.Duration(3)*time.Second)

	if err := os.RemoveAll(mntURL); err != nil {
		log.Errorf("Remove dir %s error %v", mntURL, err)
	}
}

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.RemoveAll(writeURL); err != nil {
		log.Errorf("Remove dir %s error %v", writeURL, err)
	}
}


func DeleteWorkLayer(rootURL string) {
	workURL := rootURL + "workLayer/"
	if err := os.RemoveAll(workURL); err != nil {
		log.Errorf("Remove dir %s error %v", workURL, err)
	}
}


func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}






