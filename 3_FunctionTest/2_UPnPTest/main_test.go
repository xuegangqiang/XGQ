package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"testing"

	Bytomupnp "github.com/Bytom/bytom/p2p/upnp"

	EthereumUPnP "github.com/ethereum/go-ethereum/p2p/nat"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/internetgateway1"
)

func TestBytomUPnP(t *testing.T) {
	// caps, _ := upnp.Probe()
	// info := fmt.Sprintf("Hairpin:%v,PortMapping:%v", caps.Hairpin, caps.PortMapping)
	// fmt.Println(info)
	var i, TryNum int
	var info string
	var nat Bytomupnp.NAT
	var err error

	TryNum = 100
	for i = 0; i < TryNum; i++ {
		nat, err = Bytomupnp.Discover()
		if err == nil {
			break
		}
	}
	if err != nil {
		info = fmt.Sprintf("Err:%v", err)
		fmt.Println(info)
	}

	var addr net.IP
	addr, err = nat.GetExternalAddress()
	if err != nil {
		info = fmt.Sprintf("Err:%v", err)
		fmt.Println(info)
	}
	info = fmt.Sprintf("Addr:%v", addr)
	fmt.Println(info)

	var mappedExternalPort int

	// intPort, extPort := 8001, 8001
	intPort, extPort := 48001, 58001
	// mappedExternalPort, err = nat.AddPortMapping("tcp", intPort, extPort, "20180815 upnp", 0)
	mappedExternalPort, err = nat.AddPortMapping("tcp", intPort, extPort, "Tendermint UPnP Probe", 2000)
	if err != nil {
		info = fmt.Sprintf("Err:%v", err)
		fmt.Println(info)
	}
	info = fmt.Sprintf("mappedExternalPort:%v", mappedExternalPort)
	fmt.Println(info)
}

func Test_WANPPPConnection1_GetExternalIPAddress(t *testing.T) {
	clients, errors, err := internetgateway1.NewWANPPPConnection1Clients()
	extIPClients := make([]GetExternalIPAddresser, len(clients))
	for i, client := range clients {
		extIPClients[i] = client
	}
	DisplayExternalIPResults(extIPClients, errors, err)
}

func Test_WANIPConnection_GetExternalIPAddress(t *testing.T) {
	clients, errors, err := internetgateway1.NewWANIPConnection1Clients()
	extIPClients := make([]GetExternalIPAddresser, len(clients))
	for i, client := range clients {
		extIPClients[i] = client
	}
	DisplayExternalIPResults(extIPClients, errors, err)
}

func Test_ReuseDiscoveredDevice(t *testing.T) {
	var allMaybeRootDevices []goupnp.MaybeRootDevice
	for _, urn := range []string{internetgateway1.URN_WANPPPConnection_1, internetgateway1.URN_WANIPConnection_1} {
		maybeRootDevices, err := goupnp.DiscoverDevices(internetgateway1.URN_WANPPPConnection_1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not discover %s devices: %v\n", urn, err)
		}
		allMaybeRootDevices = append(allMaybeRootDevices, maybeRootDevices...)
	}
	locations := make([]*url.URL, 0, len(allMaybeRootDevices))
	fmt.Fprintf(os.Stderr, "Found %d devices:\n", len(allMaybeRootDevices))
	for _, maybeRootDevice := range allMaybeRootDevices {
		if maybeRootDevice.Err != nil {
			fmt.Fprintln(os.Stderr, "  Failed to probe device at ", maybeRootDevice.Location.String())
		} else {
			locations = append(locations, maybeRootDevice.Location)
			fmt.Fprintln(os.Stderr, "  Successfully probed device at ", maybeRootDevice.Location.String())
		}
	}
	fmt.Fprintf(os.Stderr, "Attempt to re-acquire %d devices:\n", len(locations))
	for _, location := range locations {
		if _, err := goupnp.DeviceByURL(location); err != nil {
			fmt.Fprintf(os.Stderr, "  Failed to reacquire device at %s: %v\n", location.String(), err)
		} else {
			fmt.Fprintf(os.Stderr, "  Successfully reacquired device at %s\n", location.String())
		}
	}
}

func Test_EthereumUPnP(t *testing.T) {

	var info string
	var eupnp EthereumUPnP.Interface

	eupnp = EthereumUPnP.UPnP()

	IP, err := eupnp.ExternalIP()
	if err != nil {
		info = fmt.Sprintf("Err:%v", err)
		fmt.Println(info)
	}

	info = IP.String()
	fmt.Println(info)

	err = eupnp.AddMapping("tcp", 444, 334, "eupnp", 0)
	if err != nil {
		info = fmt.Sprintf("Err:%v", err)
		fmt.Println(info)
	}

	info = fmt.Sprintf("AddMapping succeed!")
	fmt.Println(info)

}
