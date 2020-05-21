package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	_ "path/filepath"
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





//init挂载点
func pivotRoot(root string) error{

	//syscall.Mount("", "/", "", syscall.MS_PRIVATE | syscall.MS_REC, "")

	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND | syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error:%v", err)
	}
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	if err := syscall.PivotRoot(root, pivotDir); err != nil{
		return fmt.Errorf("pivot_root %v", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err!=nil{
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	return os.Remove(pivotDir)
}


//func pivotRoot(rootfs string) error {
//	oldroot, err := syscall.Open("/", syscall.O_DIRECTORY|syscall.O_RDONLY, 0)
//	if err != nil {
//		return nil
//	}
//	defer syscall.Close(oldroot)
//	newroot, err := syscall.Open(rootfs, syscall.O_DIRECTORY|syscall.O_RDONLY, 0)
//	if err != nil {
//		return nil
//	}
//	defer syscall.Close(newroot)
//	if err := syscall.Fchdir(newroot); err != nil{
//		return err
//	}
//	if err := syscall.PivotRoot(".", "."); err != nil{
//		//if err := rootfsParentMountPrivate("."); err != nil {
//		//	return err
//		//}
//		if err := syscall.PivotRoot(".", "."); err != nil{
//			return fmt.Errorf("pivot_root %s", err)
//		}
//	}
//	if err := syscall.Fchdir(oldroot); err != nil {
//		return err
//	}
//	if err := syscall.Mount("", ".", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
//		return err
//	}
//	if err := syscall.Unmount(".", syscall.MNT_DETACH); err != nil {
//		return err
//	}
//	if err := syscall.Chdir("/"); err != nil {
//		return fmt.Errorf("chdir / %s", err)
//	}
//	return nil
//}




func setUpMount(){
	pwd, err := os.Getwd()
	if err != nil{
		log.Errorf("Get current location error %v", err)
		return
	}
	log.Infof("Current location is %s", pwd)

	pivotRoot(pwd)

	syscall.Mount("", "/", "", syscall.MS_PRIVATE | syscall.MS_REC, "")
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID | syscall.MS_STRICTATIME, "mode=755")
}


func readUserCommand() []string{
	pipe := os.NewFile(uintptr(3), "pipe")
	defer pipe.Close()
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

	setUpMount()

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