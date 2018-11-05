package main

import (
	kis "XGQ/2_KademliaSimu/KademliaInformationSystem"
	"fmt"
	"testing"

	GoCfg "github.com/Unknwon/goconfig"
)

func TestUPnPMappingAdd(t *testing.T) {
	var PortInternal, PortExternal int
	var InternalIP, ExternalIP []uint8
	var InternalIPArr [][]uint8
	var info string

	nwp := kis.NetWorkProtocol{}

	PortInternal = 123
	PortExternal = 334

	InternalIPArr = nwp.GetInternalIP()
	InternalIP = InternalIPArr[0]
	ExternalIP = nwp.GetExternalIP()

	info = fmt.Sprintf("InternalIP:%3d.%3d.%3d.%3d\nExternalIP:%3d.%3d.%3d.%3d",
		InternalIP[len(InternalIP)-4], InternalIP[len(InternalIP)-3],
		InternalIP[len(InternalIP)-2], InternalIP[len(InternalIP)-1],
		ExternalIP[len(ExternalIP)-4], ExternalIP[len(ExternalIP)-3],
		ExternalIP[len(ExternalIP)-2], ExternalIP[len(ExternalIP)-1])
	fmt.Println(info)

	if !nwp.AddUPnPMapping(PortInternal, PortExternal) {
		return
	}

	info = fmt.Sprintf("Add UPnP mapping succeed! Port[Internal:%5d,External:%5d]",
		PortInternal, PortExternal)
	fmt.Println(info)
}

func TestUPnPMappingDelete(t *testing.T) {
	var PortInternal, PortExternal int

	PortInternal = 123
	PortExternal = 334

	nwp := kis.NetWorkProtocol{}

	nwp.DeleteUPnPMapping(PortInternal, PortExternal)
}

func TestRPC(t *testing.T) {
	var cfg *GoCfg.ConfigFile
	var info, CfgFilePath, StrIP string

	CfgFilePath = ".\\cfg.ini"

	if xu.FileExist(CfgFilePath) {
		cfg, _ = GoCfg.LoadConfigFile(CfgFilePath)
	} else {
		info = fmt.Sprintf("Configure file not exist:%s", CfgFilePath)
		fmt.Println(info)
		return
	}

	StrIP, _ = cfg.GetValue("NetWorkCfg", "SeedNodeIP")
	fmt.Println(StrIP)
}

func TestSinglePeer(t *testing.T) {
	var p *kis.Peer

	p = kis.NewPeer()
	p.Initialization()
	
	p.CallRPCPing(xu.IP2Str(p.InternalIP, p.InternalPort+1))
}
