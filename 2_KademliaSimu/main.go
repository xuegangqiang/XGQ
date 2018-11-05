package main

import (
	kis "XGQ/2_KademliaSimu/KademliaInformationSystem"
	"XGQ/XgqUtility"
	"time"
)

var xu XgqUtility.Utility

func main() {
	PeerBoot()
}

func PeerBoot() {
	var p *kis.Peer
	var PingPeer bool
	var DstAddr string

	p = kis.NewPeer()

	p.Initialization()
	PingPeer = p.ExternalPort == 56067
	go p.RunRPCServer()

	if !PingPeer {
		<-p.RPCServerCh
		return
	}

	DstAddr = xu.IP2Str(p.ExternalIP, p.ExternalPort+1)
	for i := 0; i < 10; i++ {
		p.CallRPCPing(DstAddr)
		time.Sleep(time.Duration(1) * time.Second)
	}

	p.CloseRPCXClt(DstAddr)
}
