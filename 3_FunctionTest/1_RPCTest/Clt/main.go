package main

import (
	"XGQ/3_FunctionTest/1_RPCTest/Model"
	"context"
	"log"
	"time"

	"github.com/smallnest/rpcx/client"
)

func main() {

	// NAT 映射地址说明：
	// 1、Windows 自带的防火墙会阻止RPC调用(2台电脑分别使用有线和无线连接路由器)
	// 2、Windows 自带的防火墙会阻止同一台电脑上的NAT映射地址调用，改用局域网地址后OK
	// 3、360不阻止此类连接和映射

	// addr := "192.168.0.172:8972"

	// addr := "192.168.0.103:334"
	addr := "172.168.199.245:444"

	d := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	xclient := client.NewXClient("Model.ModelData", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &Model.Args{
		A: 10,
		B: 20,
	}

	for {
		reply := &Model.Reply{}
		err := xclient.Call(context.Background(), "RpcxMethod", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		time.Sleep(1e9)
	}
}
