package KademliaInformationSystem

import (
	"XGQ/1_DistributedLedgerDemo/Model"
	"context"
	"fmt"

	RPCX "github.com/smallnest/rpcx/client"
)

func (p *Peer) OpenRPCXClt(Addr string) {
	var info string
	var RPCXClt RPCX.XClient
	var sd RPCX.ServiceDiscovery

	if uint32(len(p.rpcm.RPCXCltMap)) >= p.rpcm.MaxConnectNum {
		info = fmt.Sprintf(
			"can not connect more than %d peer",
			p.rpcm.MaxConnectNum)
		fmt.Println(info)

		if KCI != nil {
			KCI.Report2Ctrl(info)
		}

		return
	}

	sd = RPCX.NewPeer2PeerDiscovery("tcp@"+Addr, "")
	RPCXClt = RPCX.NewXClient(
		"KademliaInformationSystem.RPCMethod",
		RPCX.Failtry, RPCX.RandomSelect,
		sd, RPCX.DefaultOption)

	p.rpcm.RPCXCltMap[Addr] = RPCXClt
}

func (p *Peer) CloseRPCXClt(Addr string) {
	if p.rpcm.RPCXCltMap[Addr] == nil {
		return
	}
	p.rpcm.RPCXCltMap[Addr].Close()
}

func (p *Peer) callRPCTemplate(Addr, MethodStr string, args, reply interface{}) error {
	var info string
	var err error

	if p.rpcm.RPCXCltMap[Addr] == nil {
		p.OpenRPCXClt(Addr)
	}

	err = p.rpcm.RPCXCltMap[Addr].Call(
		context.Background(), MethodStr,
		args, reply)
	if err != nil {
		info = fmt.Sprintf("%s failed: %s", MethodStr, err)
		fmt.Println(info)

		if KCI != nil {
			KCI.Report2Ctrl(info)
		}
	}

	return err
}

func (p *Peer) CallRPCPing(Addr string) bool {
	var info string
	DstPeer := new(Bucket)

	err := p.callRPCTemplate(
		Addr, "RPCPing",
		p.FormABucket(), DstPeer)

	if err != nil {
		return false
	}

	info = fmt.Sprintf(
		"%s,Reply from Svr %s",
		xu.GetCurTimeStr(),
		xu.IP2Str(DstPeer.ExternalIP, DstPeer.ExternalPort))
	fmt.Println(info)

	if KCI != nil {
		KCI.Report2Ctrl(info)
	}

	return true
}

func (p *Peer) CallRPCFindNode(Addr, UUID string) (bool, []Bucket) {
	var BucketArr []Bucket
	var args ArgsRPCFindNode
	var reply ReplyRPCFindNode

	args = ArgsRPCFindNode{}
	args.TarUUID = UUID
	args.FormFromBucket(p.FormABucket())

	err := p.callRPCTemplate(
		Addr, "RPCFindNode",
		&args, &reply)

	BucketArr = reply.FormBucketArr()

	if err != nil {
		return false, BucketArr
	}

	return true, BucketArr
}

func (p *Peer) CallRPCBlockChainSync(Addr string) {
	var args, reply *ArgsBlockChainSync
	var LocalPOW *Model.ProofOfWork
	var AllAddr []string
	var i int
	var IsBlockChainUpdated bool

	AllAddr = p.GetAllAddr()
	LocalPOW = KCI.GetPOW()
	reply = &ArgsBlockChainSync{}

	for {
		IsBlockChainUpdated = false
		args = NewArgsBlockChainSync(LocalPOW)

		for i = 0; i < len(AllAddr); i++ {
			Addr = AllAddr[i]
			if Addr == xu.IP2Str(p.ExternalIP, p.ExternalPort) {
				continue
			}

			p.callRPCTemplate(
				Addr, "RPCBlockChainSync",
				&args, &reply)

			if reply.BlockChainLen <= args.BlockChainLen {
				continue
			}

			reply.FormPOW(LocalPOW)
			LocalPOW.ClearTxPool()
			IsBlockChainUpdated = true
		}

		if !IsBlockChainUpdated {
			break
		}

		IsBlockChainUpdated = false
	}
}
