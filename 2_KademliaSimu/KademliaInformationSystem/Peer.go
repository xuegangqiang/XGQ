package KademliaInformationSystem

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	GoCfg "github.com/Unknwon/goconfig"
	uuid "github.com/satori/go.uuid"
)

type Peer struct {
	InternalIP   net.IP
	ExternalIP   net.IP
	SeedNodeIP   net.IP
	InternalPort uint32
	ExternalPort uint32
	SeedNodePort uint32
	UUID         string
	SeedNodeUUID string
	kbuckets     KBuckets

	rpcm        *RPCMethod
	RPCServerCh chan bool

	CfgFilePath string
}

func NewPeer() *Peer {
	p := Peer{}

	kb := NewKBuckets()
	p.kbuckets = *kb
	p.kbuckets.peer = &p

	u1 := uuid.Must(uuid.NewV4())
	p.UUID = fmt.Sprintf("%s", u1)

	u1 = uuid.Must(uuid.NewV4())
	p.SeedNodeUUID = fmt.Sprintf("%s", u1)

	p.InternalIP = make([]byte, 4, 4)
	for i, _ := range p.InternalIP {
		p.InternalIP[i] = byte(uint8(rand.Intn(256)))
	}
	p.InternalPort = uint32(rand.Intn(99999))
	p.ExternalIP = make([]byte, 4, 4)
	p.ExternalPort = 0

	p.rpcm = NewRPCMethod(&p)
	p.RPCServerCh = make(chan bool)

	p.CfgFilePath = ".\\Data\\NetworkCfg.ini"

	return &p
}

func (p *Peer) Initialization() {
	var info string
	var nwp NetWorkProtocol

	nwp = NetWorkProtocol{}

	p.GetIPAddr()
	p.LoadNetWorkCfg("")

	if !nwp.AddUPnPMapping(int(p.InternalPort), int(p.ExternalPort)) {
		info = fmt.Sprintf("Add UPnP mapping failed")
		fmt.Println(info)
	} else {
		info = fmt.Sprintf("Add UPnP mapping succeed! Port[Internal:%5d,External:%5d]",
			int(p.InternalPort), int(p.ExternalPort))
		fmt.Println(info)
	}

	p.SaveCfgFile("")
}

func (p *Peer) GetKBuckets() *KBuckets {
	return &p.kbuckets
}

func (p *Peer) GetIPAddr() {
	var InternalIP, ExternalIP []uint8
	var InternalIPArr [][]uint8
	var nwp NetWorkProtocol

	nwp = NetWorkProtocol{}

	InternalIPArr = nwp.GetInternalIP()
	InternalIP = InternalIPArr[0]
	ExternalIP = nwp.GetExternalIP()
	InternalIP = InternalIP[len(InternalIP)-4 : len(InternalIP)]
	if len(ExternalIP) >= 4 {
		ExternalIP = ExternalIP[len(ExternalIP)-4 : len(ExternalIP)]
	} else {
		ExternalIP = make([]byte, 4, 4)
	}
	p.InternalIP = InternalIP
	p.ExternalIP = ExternalIP
}

func (p *Peer) LoadNetWorkCfg(CfgFilePath string) {
	var nwp NetWorkProtocol
	var cfg *GoCfg.ConfigFile
	var InternalAddrStr, ExternalAddrStr, SeedNodeAddrStr string
	var UUIDStr, SeedNodeUUIDStr string
	var CfgInternalIP, CfgExternalIP []uint8
	var CfgInternalPort, CfgExternalPort uint32

	if CfgFilePath == "" {
		CfgFilePath = p.CfgFilePath
	}

	nwp = NetWorkProtocol{}

	if !xu.FileExist(CfgFilePath) {
		FileID, _ := os.Create(CfgFilePath)
		defer FileID.Close()
	}
	cfg, _ = GoCfg.LoadConfigFile(CfgFilePath)

	SeedNodeAddrStr, _ = cfg.GetValue("NetWorkCfg", "SeedNodeAddr")
	if SeedNodeAddrStr == "" {
		SeedNodeAddrStr = "192.168.0.106:112"
	}
	p.SeedNodeIP, p.SeedNodePort = xu.Str2IP(SeedNodeAddrStr)

	InternalAddrStr, _ = cfg.GetValue("NetWorkCfg", "LocalInternalAddr")
	ExternalAddrStr, _ = cfg.GetValue("NetWorkCfg", "LocalExternalAddr")
	p.ExternalPort = uint32(rand.Intn(99999))
	if ExternalAddrStr != "" {
		CfgExternalIP, CfgExternalPort = xu.Str2IP(ExternalAddrStr)
		CfgInternalIP, CfgInternalPort = xu.Str2IP(InternalAddrStr)

		if CfgInternalIP[0] != p.InternalIP[0] ||
			CfgInternalIP[1] != p.InternalIP[1] ||
			CfgInternalIP[2] != p.InternalIP[2] ||
			CfgInternalIP[3] != p.InternalIP[3] ||
			CfgExternalIP[0] != p.ExternalIP[0] ||
			CfgExternalIP[1] != p.ExternalIP[1] ||
			CfgExternalIP[2] != p.ExternalIP[2] ||
			CfgExternalIP[3] != p.ExternalIP[3] {
			nwp.DeleteUPnPMapping(int(CfgInternalPort), int(CfgExternalPort))
		} else {
			p.ExternalPort = CfgExternalPort
			p.InternalPort = CfgInternalPort
		}
	}

	UUIDStr, _ = cfg.GetValue("NetWorkCfg", "LocalUUID")
	SeedNodeUUIDStr, _ = cfg.GetValue("NetWorkCfg", "SeedNodeUUID")

	if UUIDStr != "" {
		p.UUID = UUIDStr
	}

	if SeedNodeUUIDStr != "" {
		p.SeedNodeUUID = SeedNodeUUIDStr
	}
}

func (p *Peer) SaveCfgFile(
	CfgFilePath string) (err error) {
	var cfg *GoCfg.ConfigFile
	var info string

	if CfgFilePath == "" {
		CfgFilePath = p.CfgFilePath
	}

	if !xu.FileExist(CfgFilePath) {
		FileID, _ := os.Create(CfgFilePath)
		defer FileID.Close()
	}
	cfg, _ = GoCfg.LoadConfigFile(CfgFilePath)

	cfg.SetValue("NetWorkCfg", "SeedNodeUUID", p.SeedNodeUUID)
	cfg.SetValue("NetWorkCfg", "SeedNodeAddr",
		xu.IP2Str(p.SeedNodeIP, p.SeedNodePort))
	cfg.SetSectionComments("NetWorkCfg", "# NetWorkCfg")

	cfg.SetValue("NetWorkCfg", "LocalUUID", p.UUID)
	cfg.SetValue("NetWorkCfg", "LocalInternalAddr",
		xu.IP2Str(p.InternalIP, p.InternalPort))
	cfg.SetValue("NetWorkCfg", "LocalExternalAddr",
		xu.IP2Str(p.ExternalIP, p.ExternalPort))

	err = GoCfg.SaveConfigFile(cfg, CfgFilePath)
	if err != nil {
		info = fmt.Sprintf("%s", err)
		fmt.Println(info)
	}

	return err
}

func (p *Peer) JoinNetWork(RefPeerArr *[]Peer) {
	for i, _ := range *RefPeerArr {
		p.kbuckets.AddBucket(*(*RefPeerArr)[i].FormABucket())
		p.kbuckets.NodeLookUp(p.UUID)
	}
}

func (p *Peer) FormABucket() *Bucket {
	b := Bucket{}

	b.UUID = p.UUID
	b.InternalIP = p.InternalIP
	b.ExternalIP = p.ExternalIP
	b.InternalPort = p.InternalPort
	b.ExternalPort = p.ExternalPort
	b.Timestamp = time.Now().Unix()

	return &b
}

func (p *Peer) PingAddr(
	ExternalIP net.IP, ExternalPort uint32) bool {
	return p.CallRPCPing(xu.IP2Str(ExternalIP, ExternalPort))
}

func (p *Peer) GetAllAddr() []string {
	return p.kbuckets.GetAllAddr()
}
