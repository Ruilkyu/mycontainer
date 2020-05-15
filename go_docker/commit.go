package main

import (
	"fmt"
	"os/exec"
	log "github.com/sirupsen/logrus"
)

func commitContainer(imageName string){
	mntURL := "/root/merged"
	imageURL := "/root/" + imageName + ".tar"
	fmt.Printf("%s", imageURL)
	if _,err := exec.Command("tar","-zcvf",imageURL,"-C", mntURL, ".").CombinedOutput(); err != nil{
		log.Errorf("Tar folder %s error %v", mntURL, err)
	}
}
