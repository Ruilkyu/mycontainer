package subsystems


import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)


type MemorySubSystem struct {

}


//返回cgroup名字
func (s *MemorySubSystem) Name() string{
	return "memory"
}


// 设置cgroupPath对应的cgroup的内存限制
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error{
	if subsysCgroupPath,err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != ""{
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}


//将一个进程加入到cgroupPath对应的cgroup中
func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error{
	if subsysCgroupPath,err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath,"tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	}else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}



//删除cgroupPath对应的cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath,err := GetCgroupPath(s.Name(), cgroupPath, false);err == nil{
		return os.Remove(subsysCgroupPath)
	} else {
		return err
	}
}


















