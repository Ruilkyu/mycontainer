package container

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

// 创建一个Overlay系统作为容器的根目录
func NewWorkSpace(volume string, imageName string, containerName string){
	CreateReadOnlyLayer(imageName)
	CreateWriteLayer(containerName)
	CreateWorkLayer(containerName)
	CreateMountPoint(containerName, imageName)

	// 判断是否执行挂载数据卷操作
	if volume != ""{
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(volumeURLs, containerName)
			log.Infof("NewWorkSpace volume urls %q", volumeURLs)
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

// 创建只读层
func CreateReadOnlyLayer(imageName string) error{
	unTarFolderUrl := fmt.Sprintf(ReadLayerUrl, imageName)
	imageUrl := RootUrl + "/" + imageName + ".tar"
	exist, err := PathExists(unTarFolderUrl)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", unTarFolderUrl, err)
		return err
	}
	if exist == false {
		if err := os.MkdirAll(unTarFolderUrl, 0622); err != nil {
			log.Errorf("Mkdir read layer dir %s error. %v", unTarFolderUrl, err)
			return err
		}
		if _, err := exec.Command("tar", "-xvf", imageUrl, "-C", unTarFolderUrl).CombinedOutput(); err != nil {
			log.Errorf("Untar read dir %s error %v", unTarFolderUrl, err)
			return err
		}
	}
	return nil
}


// 创建读写层
func CreateWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		log.Errorf("Mkdir write layer dir %s error. %v", writeURL, err)
	}
}

// 创建工作层
func CreateWorkLayer(containerName string) {
	workURL := fmt.Sprintf(WorkLayerUrl, containerName)
	if err := os.MkdirAll(workURL, 0777); err != nil {
		log.Errorf("Mkdir work layer dir %s error. %v", workURL, err)
	}
}


// 挂载数据卷
func MountVolume(volumeURLs []string, containerName string) error{
	// 在宿主机创建宿主目录
	parentUrl := volumeURLs[0]
	if err := os.MkdirAll(parentUrl, 0777); err != nil {
		log.Infof("Mkdir parent dir %s error. %v", parentUrl, err)
	}
	CreateVolumePoint(parentUrl)

	// 在容器里创建挂载目录
	containerUrl := volumeURLs[1]
	mntURL := fmt.Sprintf(MntUrl, containerName)
	containerVolumeURL := mntURL + "/" + containerUrl
	if err := os.MkdirAll(containerVolumeURL, 0777); err != nil {
		log.Infof("Mkdir container dir %s error. %v", containerVolumeURL, err)
	}

	// 把宿主机文件目录挂载到容器挂载点
	lower := "lowerdir=" + parentUrl + "/readLayer/"
	upper := "upperdir=" + parentUrl + "/writeLayer/"
	work := "workdir=" + parentUrl + "/workLayer/"
	parm := lower + "," + upper + "," + work
	//cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	_,err := exec.Command("mount", "-t", "overlay", "overlay", "-o", parm, containerVolumeURL).CombinedOutput()
	if err != nil {
		log.Errorf("Mount volume failed. %v", err)
		return err
	}
	return nil
}


// 创建宿主机挂载卷的read/write.work目录
func CreateVolumePoint(parentUrl string) {
	readURL := parentUrl + "/readLayer/"
	writeURL := parentUrl + "/writeLayer/"
	workURL := parentUrl + "/workLayer/"

	if err := os.MkdirAll(readURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", readURL, err)
	}

	if err := os.MkdirAll(writeURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeURL, err)
	}

	if err := os.MkdirAll(workURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", workURL, err)
	}
}



// 创建容器目录
func CreateMountPoint(containerName string, imageName string) error {
	mntUrl := fmt.Sprintf(MntUrl, containerName)
	if err := os.MkdirAll(mntUrl, 0777); err != nil {
		log.Errorf("Mkdir mountpoint dir %s error. %v", mntUrl, err)
	}
	tmpReadLayer := fmt.Sprintf(ReadLayerUrl, imageName)
	tmpWriteLayer := fmt.Sprintf(WriteLayerUrl, containerName)
	tmpwork := fmt.Sprintf(WorkLayerUrl, containerName)
	lower := "lowerdir=" + tmpReadLayer
	upper := "upperdir=" + tmpWriteLayer
	work := "workdir=" + tmpwork
	parm := lower + "," + upper + "," + work
	//cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	_, err := exec.Command("mount", "-t", "overlay", "overlay", "-o", parm, mntUrl).CombinedOutput()
	if err != nil {
		log.Errorf("Run command for creating mount point failed %v", err)
		return err
	}
	return nil
}

// 当容器退出时删除Overlay系统
func DeleteWorkSpace(volume string, containerName string){
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if (length == 2 && volumeURLs[0] != "" && volumeURLs[1] != ""){
			DeleteMountPointWithVolume(volumeURLs, containerName)
		} else {
			DeleteMountPoint(containerName)
		}
	} else {
		DeleteMountPoint(containerName)
	}
	DeleteWriteLayer(containerName)
	DeleteWorkLayer(containerName)
}


func DeleteMountPoint(containerName string) error{
	mntURL := fmt.Sprintf(MntUrl, containerName)
	_, err := exec.Command("umount", mntURL).CombinedOutput()
	if err != nil{
		log.Errorf("Unmount %s error %v", mntURL, err)
		return err
	}

	if err := os.RemoveAll(mntURL); err != nil {
		log.Errorf("Remove mountpoint dir %s error %v", mntURL, err)
		return err
	}
	return nil
}


func DeleteMountPointWithVolume(volumeURLs []string, containerName string) error{
	// 卸载容器里的volume挂载点
	mntURL := fmt.Sprintf(MntUrl, containerName)
	containerUrl := mntURL + "/" + volumeURLs[1]
	_,err := exec.Command("umount", containerUrl).CombinedOutput()
	if err != nil {
		log.Errorf("Umount volume %s failed. %v", containerUrl, err)
	}

	// 卸载整个容器挂载点
	if err != nil{
		log.Errorf("Unmount %s error %v", mntURL, err)
		return err
	}

	if err := os.RemoveAll(mntURL); err != nil {
		log.Errorf("Remove mountpoint dir %s error %v", mntURL, err)
		return err
	}
	return nil
}


func DeleteWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.RemoveAll(writeURL); err != nil {
		log.Errorf("Remove writeLayer dir %s error %v", writeURL, err)
	}
}


func DeleteWorkLayer(containerName string) {
	workURL := fmt.Sprintf(WorkLayerUrl, containerName)
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

