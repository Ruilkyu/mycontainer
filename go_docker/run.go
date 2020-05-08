package main


import (
	log "github.com/sirupsen/logrus"
	"go_docker/container"
	"go_docker/cgroups"
	"go_docker/cgroups/subsystems"
	"os"
	"strings"
)


//func Run(tty bool, command string){
//	parent := container.NewParentProcess(tty, command)
//	if err := parent.Start(); err != nil {
//		log.Error(err)
//	}
//	parent.Wait()
//	os.Exit(-1)
//}



func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}


	 cgroupManager := cgroups.NewCgroupManager("go_docker")
	 defer cgroupManager.Destroy()
	 cgroupManager.Set(res)
	 cgroupManager.Apply(parent.Process.Pid)

	 sendInitCommand(comArray, writePipe)

	 parent.Wait()


	 mntURL := "/root/merged/"
	 rootURL := "/root/"
	 container.DeleteWorkSpace(rootURL, mntURL)

	 os.Exit(0)
}


func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
