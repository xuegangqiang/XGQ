package main

import (
	"XGQ/3_FunctionTest/1_RPCTest/Model"

	"github.com/smallnest/rpcx/server"
)

func main() {

	s := server.NewServer()
	s.RegisterName("Model.ModelData", new(Model.ModelData), "")
	// s.Serve("tcp", ":8972")

	// s.Serve("tcp", "172.168.199.245:444")
	// s.Serve("tcp", "192.168.0.116:334")

	s.Serve("tcp", ":334")
}
