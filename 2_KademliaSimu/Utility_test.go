package main

import (
	"fmt"
	"os"
	"testing"

	GoCfg "github.com/Unknwon/goconfig"
)

func TestGoCfg(t *testing.T) {
	var cfg *GoCfg.ConfigFile
	var err error
	var info, CfgFilePath string

	CfgFilePath = ".\\cfg.ini"

	if !xu.FileExist(CfgFilePath) {
		FileID, _ := os.Create(CfgFilePath)
		defer FileID.Close()
	}
	cfg, _ = GoCfg.LoadConfigFile(CfgFilePath)

	cfg.SetValue("NetWorkCfg", "SeedNodeIP", "192.168.0.113")
	cfg.SetSectionComments("NetWorkCfg", "# NetWorkCfg")

	err = GoCfg.SaveConfigFile(cfg, CfgFilePath)
	if err != nil {
		info = fmt.Sprintf("%s", err)
		fmt.Println(info)
	}

}

func TestMap(t *testing.T) {
	var m map[string]bool
	var b bool
	var c int

	// p := *(kis.NewPeer())
	// p.Initialization()

	m = make(map[string]bool, 1)
	m["aa"] = true
	c = len(m)
	m["aa1"] = true
	c = len(m)
	m["aa3"] = true
	c = len(m)
	m["aa4"] = true
	c = len(m)
	m["aa5"] = true
	c = len(m)

	delete(m, "aa")
	c = len(m)

	b = m["aa5"]

	if b {
		m["aa1"] = c == 1
	}

}

func TestSlice(t *testing.T) {
	var a []uint8
	a = make([]uint8, 5, 8)
	info := fmt.Sprintf("%d,%d", len(a), cap(a))
	fmt.Println(info)
}
