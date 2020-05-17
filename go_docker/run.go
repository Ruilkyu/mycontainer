package main


import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"go_docker/container"
	"go_docker/cgroups"
	"go_docker/cgroups/subsystems"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)


//func Run(tty bool, command string){
//	parent := container.NewParentProcess(tty, command)
//	if err := parent.Start(); err != nil {
//		log.Error(err)
//	}
//	parent.Wait()
//	os.Exit(-1)
//}



func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume string, containerName string, logfile bool) {
	parent, writePipe := container.NewParentProcess(tty, volume, containerName,logfile)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//record container info
	containerName, err := recordContainerInfo(parent.Process.Pid, comArray, containerName)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	 cgroupManager := cgroups.NewCgroupManager("go_docker")
	 defer cgroupManager.Destroy()
	 cgroupManager.Set(res)
	 cgroupManager.Apply(parent.Process.Pid)

	 sendInitCommand(comArray, writePipe)

	 if tty{
		 parent.Wait()
		 mntURL := "/root/merged/"
		 rootURL := "/root/"
		 container.DeleteWorkSpace(rootURL, mntURL, volume)
		 deleteContainerInfo(containerName)
		 os.Exit(0)
	 }
}


func deleteContainerInfo(containerId string){
	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerId)
	if err := os.RemoveAll(dirUrl); err != nil {
		log.Errorf("Remove dir %s error %v", dirUrl, err)
	}
}


func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}


func randStringBytes(n int) string{
	letterBytes := "0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b{
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func recordContainerInfo(containerPID int, commandArray []string, containerName string)(string,error){
	id := randStringBytes(10)
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")

	if containerName == ""{
		containerName = id
	}
	containerInfo := &container.ContainerInfo{
		Id: id,
		Pid: strconv.Itoa(containerPID),
		Command: command,
		CreatedTime: createTime,
		Status: container.RUNNING,
		Name: containerName,
	}

	//容器信息对象的json序列化为字符串
	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil{
		log.Errorf("Record container info error %v", err)
		return "",err
	}
	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation,containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirUrl, err)
		return "",err
	}
	fileName := dirUrl + "/" + container.ConfigName
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error %v", fileName, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return "", err
	}
	return containerName, nil
}
