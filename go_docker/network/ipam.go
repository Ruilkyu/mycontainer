package network

import (
	"encoding/json"
	"net"
	"os"
	"path"
	"strings"
	log "github.com/sirupsen/logrus"
)

//存放IP地址分配信息
const ipamDefaultAllocatorPath = "/var/run/mycontainer/network/ipam/subnet.json"
type IPAM struct{
	//分配文件存放位置
	SubnetAllocatorPath string
	// 网段和位图的数组，key是网段,value是分配的位图数组
	Subnets *map[string]string
}

// 初始化一个IPAM对象
var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

// 加载网段分配信息
func (ipam *IPAM) load() error {
	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	// 打开并读取存储文件
	subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
	defer subnetConfigFile.Close()
	if err != nil {
		return err
	}
	subnetJson := make([]byte, 2000)
	n, err := subnetConfigFile.Read(subnetJson)
	if err != nil {
		return err
	}
    // 将文件的内容反序列化出IP的分配信息
	err = json.Unmarshal(subnetJson[:n], ipam.Subnets)
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
		return err
	}
	return nil
}

// 存储网段分配信息
func (ipam *IPAM) dump() error {
	// 检查存储文件所在的文件夹是否存在，不存在创建，path.Split分割目录和文件
	ipamConfigFileDir, _ := path.Split(ipam.SubnetAllocatorPath)
	if _, err := os.Stat(ipamConfigFileDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(ipamConfigFileDir, 0644)
		} else {
			return err
		}
	}
	// 打开存储文件
	subnetConfigFile, err := os.OpenFile(ipam.SubnetAllocatorPath, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0644)
	defer subnetConfigFile.Close()
	if err != nil {
		return err
	}
    // 序列化ipam对象到json
	ipamConfigJson, err := json.Marshal(ipam.Subnets)
	if err != nil {
		return err
	}
    // 将序列化的json写入配置文件
	_, err = subnetConfigFile.Write(ipamConfigJson)
	if err != nil {
		return err
	}
	return nil
}

// 在网段中分配一个可用的IP地址
func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	// 存放网段中地址分配信息的数组
	ipam.Subnets = &map[string]string{}

	// 从文件中加载已经分配的网段信息
	err = ipam.load()
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
	}
    // 解析字符串
	_, subnet, _ = net.ParseCIDR(subnet.String())
    // 获取IP掩码的位数和固定位数
	one, size := subnet.Mask.Size()
	// 如果之前没有分配这个网段，则初始化网段的分配配置
	if _, exist := (*ipam.Subnets)[subnet.String()]; !exist {
		(*ipam.Subnets)[subnet.String()] = strings.Repeat("0", 1 << uint8(size - one))
	}
    // 遍历网段的位图数组
	for c := range((*ipam.Subnets)[subnet.String()]) {
		if (*ipam.Subnets)[subnet.String()][c] == '0' {
			ipalloc := []byte((*ipam.Subnets)[subnet.String()])
			ipalloc[c] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc)
			ip = subnet.IP
			for t := uint(4); t > 0; t-=1 {
				[]byte(ip)[4-t] += uint8(c >> ((t - 1) * 8))
			}
			// 由于IP地址是从1开始分配，所以最后再加1
			ip[3]+=1
			break
		}
	}
    //将分配的结果保存到文件中
	ipam.dump()
	return
}

// 地址释放
func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = &map[string]string{}
	// 存放网段中地址分配信息的数组
	_, subnet, _ = net.ParseCIDR(subnet.String())
    // 从文件中加载网段的分配信息
	err := ipam.load()
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
	}
    // 计算ip地址在网段位图数组索引位置
	c := 0
	// 将ip地址转换成4个字节的表示方式
	releaseIP := ipaddr.To4()
	// 由于ip地址是从1开始,所以转换成索引应减1
	releaseIP[3]-=1
	// 释放IP，获得索引的方式是IP地址的每一位相减之后分别左移，将对应的数值加到索引上
	for t := uint(4); t > 0; t-=1 {
		c += int(releaseIP[t-1] - subnet.IP[t-1]) << ((4-t) * 8)
	}
    // 将分配的位图数组中索引位置的值置为0
	ipalloc := []byte((*ipam.Subnets)[subnet.String()])
	ipalloc[c] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)
    // 保存释放掉IP之后的网段分配信息
	ipam.dump()
	return nil
}



