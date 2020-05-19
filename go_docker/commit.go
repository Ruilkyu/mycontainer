package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"go_docker/container"
	"os/exec"
)

func commitContainer(containerName string,imageName string){
	mntURL := fmt.Sprintf(container.MntUrl, containerName)
	mntURL += "/"
	imageURL := container.RootUrl + "/" + imageName + ".tar"

	if _,err := exec.Command("tar","-zcvf",imageURL,"-C", mntURL, ".").CombinedOutput(); err != nil{
		log.Errorf("Tar folder %s error %v", mntURL, err)
	}
}
