package main

// To block cmd window:
// go build -ldflags="-H windowsgui"

// You can instruct the linker to strip debug symbols by using:
// go install -ldflags '-s'

import (
	"XGQ/1_DistributedLedgerDemo/Controller"
)

func main() {
	var c Controller.Controller
	c = Controller.Controller{}
	c.Run()
}
