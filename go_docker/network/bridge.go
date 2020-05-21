package network


import (
	"fmt"
	"net"
	"strings"
	"time"
	"github.com/vishvananda/netlink"
	log "github.com/sirupsen/logrus"
	"os/exec"
)


type BridgeNetworkDriver struct {

}

func (d *BridgeNetworkDriver) Name() string{
	return "bridge"
}

// 配置bridge网络
func (d *BridgeNetworkDriver) Create(subnet string, name string) (*Network, error) {
	// 获取网关IP和网段
	ip, ipRange, _ := net.ParseCIDR(subnet)
	ipRange.IP = ip
	// 初始化网络对象
	n := &Network {
		Name: name,
		IpRange: ipRange,
		Driver: d.Name(),
	}
	err := d.initBridge(n)
	if err != nil {
		log.Errorf("error init bridge: %v", err)
	}
    // 返回配置好的网络
	return n, err
}

// 初始化网桥(Bridge设备)
func (d *BridgeNetworkDriver) initBridge(n *Network) error {
	// 创建bridge虚拟设备
	bridgeName := n.Name
	if err := createBridgeInterface(bridgeName); err != nil {
		return fmt.Errorf("Error add bridge： %s, Error: %v", bridgeName, err)
	}
	// 设置bridge设备的地址和路由
	gatewayIP := *n.IpRange
	gatewayIP.IP = n.IpRange.IP

	if err := setInterfaceIP(bridgeName, gatewayIP.String()); err != nil {
		return fmt.Errorf("Error assigning address: %s on bridge: %s with an error of: %v", gatewayIP, bridgeName, err)
	}
    // 启动bridge设备
	if err := setInterfaceUP(bridgeName); err != nil {
		return fmt.Errorf("Error set bridge up: %s, Error: %v", bridgeName, err)
	}
	// 设置iptables的snat规则
	if err := setupIPTables(bridgeName, n.IpRange); err != nil {
		return fmt.Errorf("Error setting iptables for %s: %v", bridgeName, err)
	}

	return nil
}

// 创建bridge虚拟设备
func createBridgeInterface(bridgeName string) error {
	// 检查是否存在同名bridge设备
	_, err := net.InterfaceByName(bridgeName)
	if err == nil || !strings.Contains(err.Error(), "no such network interface") {
		return err
	}
	// 初始化一个netlink的Link基础对象，Link的名字即bridge虚拟设备的名字
	la := netlink.NewLinkAttrs()
	la.Name = bridgeName
    // 创建netlink的bridge对象
	br := &netlink.Bridge{LinkAttrs: la}
	// 创建bridge虚拟网络设备
	if err := netlink.LinkAdd(br); err != nil {
		return fmt.Errorf("Bridge creation failed for bridge %s: %v", bridgeName, err)
	}
	return nil
}

// 设置bridge设备的地址和路由
func setInterfaceIP(name string, rawIP string) error {
	retries := 2
	var iface netlink.Link
	var err error
	for i := 0; i < retries; i++ {
		iface, err = netlink.LinkByName(name)
		if err == nil {
			break
		}
		log.Debugf("error retrieving new bridge netlink link [ %s ]... retrying", name)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("Abandoning retrieving the new bridge link from netlink, Run [ ip link ] to troubleshoot the error: %v", err)
	}
	// netlink.ParseIPNet是对net.ParseCIDR的一个封装，整合了返回值中的IP和net
	//ipNet, err := netlink.ParseIPNet(rawIP)
	addr, _ := netlink.ParseAddr(rawIP)
	if err != nil {
		return err
	}
	// 给网络接口配置地址路由
	//addr := &netlink.Addr{ipNet, "", 0, 0, nil}
	return netlink.AddrAdd(iface, addr)
}

// 启动bridge设备
func setInterfaceUP(interfaceName string) error {
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return fmt.Errorf("Error retrieving a link named [ %s ]: %v", iface.Attrs().Name, err)
	}
	// 设置接口状态为UP状态
	if err := netlink.LinkSetUp(iface); err != nil {
		return fmt.Errorf("Error enabling interface for %s: %v", interfaceName, err)
	}
	return nil
}

// 设置iptables的snat规则
func setupIPTables(bridgeName string, subnet *net.IPNet) error {
	// 创建iptables命令
	iptablesCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
	//执行iptables命令配置SNAT规则
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("iptables Output, %v", output)
	}
	return err
}

// 删除bridge网络
func (d *BridgeNetworkDriver) Delete(network Network) error {
	bridgeName := network.Name
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}
	// 删除网络对应的linux bridge设备
	return netlink.LinkDel(br)
}


// 连接linux网桥到容器网络端点
func (d *BridgeNetworkDriver) Connect(network *Network, endpoint *Endpoint) error {
	// 获取网络名，即linux bridge的接口名
	bridgeName := network.Name
	// 通过接口名获取linux bridge接口的对象和接口属性
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}
    // 创建veth接口的配置
	la := netlink.NewLinkAttrs()
	// 由于linux接口名限制，名字取endpoints ID前5位
	la.Name = endpoint.ID[:5]
	// 通过设置veth接口的master属性，设置这个veth的一端挂载到linux bridge上
	la.MasterIndex = br.Attrs().Index
    // 配置veth属性另外一端的名字cif-{endpoint ID 前5位}
	endpoint.Device = netlink.Veth{
		LinkAttrs: la,
		PeerName:  "cif-" + endpoint.ID[:5],
	}
    // 创建veth接口
	if err = netlink.LinkAdd(&endpoint.Device); err != nil {
		return fmt.Errorf("Error Add Endpoint Device: %v", err)
	}
    // 设置veth接口启动
	if err = netlink.LinkSetUp(&endpoint.Device); err != nil {
		return fmt.Errorf("Error Add Endpoint Device: %v", err)
	}
	return nil
}


func (d *BridgeNetworkDriver) Disconnect(network Network, endpoint *Endpoint) error {
	return nil
}
