package network

import "net"

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