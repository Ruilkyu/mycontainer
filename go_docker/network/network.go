package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"runtime"
	"strings"
	"text/tabwriter"
	"go_docker/container"
)

var (
	defaultNetworkPath = "/var/run/mycontainer/network/network"
	drivers = map[string]NetworkDriver{}
	networks = map[string]*Network{}
)


//网络
type Network struct{
	Name string          //网络名
	IpRange *net.IPNet   //地址段
	Driver string        //网络驱动名
}

//网络端点
type Endpoint struct{
	ID          string `json:"id"`
	Device      netlink.Veth `json:"dev"`
	IPAddress   net.IP `json:"ip"`
	MacAddress  net.HardwareAddr `json:"mac"`
	Network     *Network
	PortMapping []string
}

//网络驱动
type NetworkDriver interface{
	Name() string                                         //驱动名
	Create(subnet string, name string) (*Network, error)  //创建网络
	Delete(network Network) error                         //删除网络
	Connect(network *Network, endpoint *Endpoint) error   //连接容器网络端点到网络
	Disconnect(network Network, endpoint *Endpoint) error //从网络上移除容器网络端点
}


// 创建网络
func CreateNetwork(driver, subnet, name string) error {
	//将字符串解析成net.IPNet
	_, cidr, _ := net.ParseCIDR(subnet)
	//通过IPMA分配网关IP
	gatewayIp, err := ipAllocator.Allocate(cidr)
	if err != nil {
		return err
	}
	cidr.IP = gatewayIp
    //通过指定的drive驱动创建网络
	nw, err := drivers[driver].Create(cidr.String(), name)
	if err != nil {
		return err
	}
    //保存网络信息，将网络信息保存在文件中，便于查询和网络连接网络端点
	return nw.dump(defaultNetworkPath)
}


//将网络对象保存到指定目录文件
func (nw *Network) dump(dumpPath string) error {
	//判断目录是否存在，不存在则创建
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dumpPath, 0644)
		} else {
			return err
		}
	}
    //用网络名作为文件名
	nwPath := path.Join(dumpPath, nw.Name)
	// 打开保存到文件，参数：存在内容则清空、只写入、不存在则创建
	nwFile, err := os.OpenFile(nwPath, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("error：", err)
		return err
	}
	defer nwFile.Close()
    // 序列化网络对象到json的字符串
	nwJson, err := json.Marshal(nw)
	if err != nil {
		log.Errorf("error：", err)
		return err
	}
    // 将网络对象的字符串写入到文件中
	_, err = nwFile.Write(nwJson)
	if err != nil {
		log.Errorf("error：", err)
		return err
	}
	return nil
}

//从指定目录文件读取网络对象信息
func (nw *Network) load(dumpPath string) error {
	// 打开配置信息
	nwConfigFile, err := os.Open(dumpPath)
	defer nwConfigFile.Close()
	if err != nil {
		return err
	}
	// 从配置文件读取网络的json配置
	nwJson := make([]byte, 2000)
	n, err := nwConfigFile.Read(nwJson)
	if err != nil {
		return err
	}
    //通过json字符串反序列化出网络对象
	err = json.Unmarshal(nwJson[:n], nw)
	if err != nil {
		log.Errorf("Error load nw info", err)
		return err
	}
	return nil
}


// 连接到指定网络
func Connect(networkName string, cinfo *container.ContainerInfo) error {
	//从networks字典读取容器连接到的网络信息
	network, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("No Such Network: %s", networkName)
	}
	// 分配容器IP地址
	ip, err := ipAllocator.Allocate(network.IpRange)
	if err != nil {
		return err
	}
	// 创建网络端点
	ep := &Endpoint{
		ID: fmt.Sprintf("%s-%s", cinfo.Id, networkName),
		IPAddress: ip,
		Network: network,
		PortMapping: cinfo.PortMapping,
	}
	// 调用网络对应的网络驱动挂载和配置网络端点
	if err = drivers[network.Driver].Connect(network, ep); err != nil {
		return err
	}
	// 到容器的namespace配置容器网络、设备IP地址和路由
	if err = configEndpointIpAddressAndRoute(ep, cinfo); err != nil {
		return err
	}
    //配置容器到宿主机的端口映射
	return configPortMapping(ep, cinfo)
}

// 离开指定网络
func Disconnect(networkName string, cinfo *container.ContainerInfo) error {
	return nil
}


//加载网络配置目录所有网络配置信息到networks字典中
func Init() error {
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver
    // 判断网络的配置目录是否存在，不存在则创建
	if _, err := os.Stat(defaultNetworkPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(defaultNetworkPath, 0644)
		} else {
			return err
		}
	}
    // 检查网络配置目录中所有文件，并通过第二个函数参数处理每一个文件
	filepath.Walk(defaultNetworkPath, func(nwPath string, info os.FileInfo, err error) error {
		//if strings.HasSuffix(nwPath, "/") {
		//	return nil
		//}
		// 如果是目录则跳过
		if info.IsDir(){
			return nil
		}
		//加载文件名作为网络名
		_, nwName := path.Split(nwPath)
		nw := &Network{
			Name: nwName,
		}
        // 加载网络配置信息
		if err := nw.load(nwPath); err != nil {
			log.Errorf("error load network: %s", err)
		}
        // 将网络的配置信息加入到networks字典中
		networks[nwName] = nw
		return nil
	})
	return nil
}

//展示创建的网络
func ListNetwork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tIpRange\tDriver\n")
	// 遍历网络信息
	for _, nw := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n", nw.Name, nw.IpRange.String(), nw.Driver,)
	}
	// 输出到标准输出
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}

// 删除网络
func DeleteNetwork(networkName string) error {
	nw, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("No Such Network: %s", networkName)
	}
    // 调用IPAM的实例ipAllocator释放网络网关的IP
	if err := ipAllocator.Release(nw.IpRange, &nw.IpRange.IP); err != nil {
		return fmt.Errorf("Error Remove Network gateway ip: %s", err)
	}
    // 调用网络驱动删除网络创建的设备与配置
	if err := drivers[nw.Driver].Delete(*nw); err != nil {
		return fmt.Errorf("Error Remove Network DriverError: %s", err)
	}
    // 从网络的配置目录中删除该网络对应的配置文件
	return nw.remove(defaultNetworkPath)
}


// 删除网络的配置文件
func (nw *Network) remove(dumpPath string) error {
	if _, err := os.Stat(path.Join(dumpPath, nw.Name)); err != nil {
		// 如果不存在直接返回
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	} else {
		// 移除这个网络对应的配置文件
		return os.Remove(path.Join(dumpPath, nw.Name))
	}
}

// 配置网络设备及路由
func configEndpointIpAddressAndRoute(ep *Endpoint, cinfo *container.ContainerInfo) error {
	// 获取veth另一端
	peerLink, err := netlink.LinkByName(ep.Device.PeerName)
	if err != nil {
		return fmt.Errorf("fail config endpoint: %v", err)
	}
	// 将容器的网络端点加入到容器的网络空间
	defer enterContainerNetns(&peerLink, cinfo)()
    // 获取容器的ip地址及网段，用于配置容器的内部接口
	interfaceIP := *ep.Network.IpRange
	interfaceIP.IP = ep.IPAddress
    // 设置容器内veth端点到你ip
	if err = setInterfaceIP(ep.Device.PeerName, interfaceIP.String()); err != nil {
		return fmt.Errorf("%v,%s", ep.Network, err)
	}
    // 启动容器内veth端点
	if err = setInterfaceUP(ep.Device.PeerName); err != nil {
		return err
	}
    // 启动lo网卡，开启127.0.0.1访问
	if err = setInterfaceUP("lo"); err != nil {
		return err
	}
    // 容器内的外部请求都通过容器内的veth端点访问
	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")

	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Gw: ep.Network.IpRange.IP,
		Dst: cidr,
	}
    // 添加路由到容器的网络空间
	if err = netlink.RouteAdd(defaultRoute); err != nil {
		return err
	}

	return nil
}

// 进入容器的网络命名空间（net namespace）
func enterContainerNetns(enLink *netlink.Link, cinfo *container.ContainerInfo) func() {
	// 找到容器的网络命名空间，打开/proc/[pid]/ns/net文件的文件描述符就可以操作net namespace
	f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", cinfo.Pid), os.O_RDONLY, 0)
	if err != nil {
		log.Errorf("error get container net namespace, %v", err)
	}
    // 获取文件句柄
	nsFD := f.Fd()
	// 锁定当前程序执行的线程
	runtime.LockOSThread()

	// 把网络端点veth peer另外一端移到容器的namespace中
	if err = netlink.LinkSetNsFd(*enLink, int(nsFD)); err != nil {
		log.Errorf("error set link netns , %v", err)
	}

	// 获取当前的网络namespace
	origns, err := netns.Get()
	if err != nil {
		log.Errorf("error get current netns, %v", err)
	}

	// 设置当前进程到新的网络namespace，并在函数执行完成之后再恢复到之前的namespace
	if err = netns.Set(netns.NsHandle(nsFD)); err != nil {
		log.Errorf("error set netns, %v", err)
	}
	return func () {
		// 恢复到之前获取到的之前的网络命名空间
		netns.Set(origns)
		// 关闭namespace文件
		origns.Close()
		// 取消对当前程序的线程锁定
		runtime.UnlockOSThread()
		// 关闭namespace文件
		f.Close()
	}
}

// 配置端口映射
func configPortMapping(ep *Endpoint, cinfo *container.ContainerInfo) error {
	// 遍历容器端口映射列表
	for _, pm := range ep.PortMapping {
		// 分成宿主机的端口和容器的端口
		portMapping :=strings.Split(pm, ":")
		if len(portMapping) != 2 {
			log.Errorf("port mapping format error, %v", pm)
			continue
		}
		// 配置iptables的prerouting中的dnat
		iptablesCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s",
			portMapping[0], ep.IPAddress.String(), portMapping[1])
		cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
		//err := cmd.Run()
		output, err := cmd.Output()
		if err != nil {
			log.Errorf("iptables Output, %v", output)
			continue
		}
	}
	return nil
}