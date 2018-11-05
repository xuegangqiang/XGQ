package KademliaInformationSystem

import (
	"fmt"
	"net"

	EthereumUPnP "github.com/ethereum/go-ethereum/p2p/nat"
)

type NetWorkProtocol struct {
}

func (nwp *NetWorkProtocol) GetInternalIP() [][]byte {
	var IP [][]byte
	var MaxIPNum, IPNum int
	var addrs []net.Addr
	var address net.Addr
	var err error

	addrs, err = net.InterfaceAddrs()
	if err != nil {
		IP = make([][]byte, 0, 0)
		return IP
	}

	MaxIPNum = len(addrs)
	IPNum = 0
	IP = make([][]byte, MaxIPNum, MaxIPNum)

	for _, address = range addrs {
		ipnet, ok := address.(*net.IPNet)

		// Filtrate loop back address
		// Only take 255.255.255.0
		if !ok ||
			ipnet.IP.IsLoopback() ||
			ipnet.IP.To4() == nil ||
			ipnet.Mask[2] != 255 {
			continue
		}

		IP[IPNum] = ipnet.IP
		IPNum++
	}

	return IP[0:IPNum]
}

// Copy go-ethereum's implement
func (nwp *NetWorkProtocol) GetExternalIP() []byte {
	var IP []byte
	var eupnp EthereumUPnP.Interface

	eupnp = EthereumUPnP.UPnP()
	IP, _ = eupnp.ExternalIP()

	return IP
}

// Copy go-ethereum's UPnP implement
func (nwp *NetWorkProtocol) AddUPnPMapping(PortInternal, PortExternal int) bool {
	var info string
	var eupnp EthereumUPnP.Interface
	var err error

	eupnp = EthereumUPnP.UPnP()

	err = eupnp.AddMapping("tcp", PortExternal, PortInternal, "eupnp", 0)
	if err != nil {
		info = fmt.Sprintf("Add UPnP mapping error:%v", err)
		fmt.Println(info)
		return false
	}

	return true
}

// Copy go-ethereum's UPnP implement
func (nwp *NetWorkProtocol) DeleteUPnPMapping(PortInternal, PortExternal int) bool {
	var info string
	var eupnp EthereumUPnP.Interface
	var err error

	eupnp = EthereumUPnP.UPnP()

	err = eupnp.DeleteMapping("tcp", PortExternal, PortInternal)
	if err != nil {
		info = fmt.Sprintf("Delete UPnP mapping error:%v", err)
		fmt.Println(info)
		return false
	}

	return true
}
