package KademliaInformationSystem

import (
	"context"
	"fmt"
	"net"

	RPCX "github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"

	"XGQ/1_DistributedLedgerDemo/Model"
)

// NAT & RPC access note:
// 1. Windows firewall will block tow device's RPC (connected to wired & wireless network)
// 2. Windows firewall will block one device's RPC with NAT address,but it's OK with local address
// 3. 360 is OK with NAT address & RPC

// RPCX data transfer note:
// 1. Test result shows that args(or reply) does not support substructures
//      such as: type Args struct { b Bucket }
//      If you do that, the substructure "b" will not be assigned
// 2. Test result shows that the first letter of the member's name must be capitalized
//      thus it has public access
//      such as: type Args struct { a uint8 }
//      If you do that, the member "a" will not be assigned
// 3. Yon can not make a new reply pointer and assign that to the original reply address
//      such as: reply = &Reply{}
//      If you do that, the reply address will be covered with a new address,
//      and the original reply address(which you really wanted) will not be assigned

type RPCMethod struct {
	p             *Peer
	MaxConnectNum uint32
	RPCXCltMap    map[string]RPCX.XClient
}

func NewRPCMethod(p *Peer) *RPCMethod {
	rpcm := &RPCMethod{}

	rpcm.p = p
	rpcm.MaxConnectNum = 10
	rpcm.RPCXCltMap = make(map[string]RPCX.XClient, rpcm.MaxConnectNum)

	return rpcm
}

func (p *Peer) RunRPCServer() {
	// Note:
	// 1. Never listen to external address(NAT device do that)
	// 2. You can only listen to internal port(you may have multiple address)
	s := server.NewServer()
	s.RegisterName("KademliaInformationSystem.RPCMethod", p.rpcm, "")

	// addr := xu.IP2Str(p.InternalIP, p.InternalPort)
	addr := fmt.Sprintf(":%d", p.InternalPort)
	s.Serve("tcp", addr)
}

func (rpcm *RPCMethod) RPCPing(
	cxt context.Context,
	args *Bucket,
	reply *Bucket) error {
	*reply = *(rpcm.p.FormABucket())

	str := fmt.Sprintf(
		"%s,Clt:%s call RPCPing from Svr:%s",
		xu.GetCurTimeStr(),
		xu.IP2Str(args.ExternalIP, args.ExternalPort),
		xu.IP2Str(reply.ExternalIP, reply.ExternalPort))
	fmt.Println(str)

	return nil
}

type ArgsRPCFindNode struct {
	Bucket
	TarUUID string
}

func (a *ArgsRPCFindNode) FormABucket() Bucket {
	var b Bucket
	b = Bucket{}

	b.UUID = a.UUID
	b.InternalIP = a.InternalIP
	b.ExternalIP = a.ExternalIP
	b.InternalPort = a.InternalPort
	b.ExternalPort = a.ExternalPort
	b.Timestamp = a.Timestamp

	return b
}

func (a *ArgsRPCFindNode) FormFromBucket(b *Bucket) {
	a.UUID = b.UUID
	a.InternalIP = b.InternalIP
	a.ExternalIP = b.ExternalIP
	a.InternalPort = b.InternalPort
	a.ExternalPort = b.ExternalPort
	a.Timestamp = b.Timestamp
}

type ReplyRPCFindNode struct {
	UUIDArr         []string
	InternalIPArr   []net.IP
	ExternalIPArr   []net.IP
	InternalPortArr []uint32
	ExternalPortArr []uint32
	TimestampArr    []int64
}

func (r *ReplyRPCFindNode) FormBucketArr() []Bucket {
	var BucketArr []Bucket
	var i, ItemNum int

	ItemNum = len(r.UUIDArr)
	BucketArr = make([]Bucket, ItemNum, ItemNum)

	for i = 0; i < ItemNum; i++ {
		BucketArr[i].UUID = r.UUIDArr[i]
		BucketArr[i].InternalIP = r.InternalIPArr[i]
		BucketArr[i].ExternalIP = r.ExternalIPArr[i]
		BucketArr[i].InternalPort = r.InternalPortArr[i]
		BucketArr[i].ExternalPort = r.ExternalPortArr[i]
		BucketArr[i].Timestamp = r.TimestampArr[i]
	}

	return BucketArr
}

func (r *ReplyRPCFindNode) FormFromBucketArr(BucketArr *[]Bucket) {
	var i, ItemNum int

	ItemNum = len(*BucketArr)
	r.UUIDArr = make([]string, ItemNum, ItemNum)
	r.InternalIPArr = make([]net.IP, ItemNum, ItemNum)
	r.ExternalIPArr = make([]net.IP, ItemNum, ItemNum)
	r.InternalPortArr = make([]uint32, ItemNum, ItemNum)
	r.ExternalPortArr = make([]uint32, ItemNum, ItemNum)
	r.TimestampArr = make([]int64, ItemNum, ItemNum)

	for i = 0; i < ItemNum; i++ {
		r.UUIDArr[i] = (*BucketArr)[i].UUID
		r.InternalIPArr[i] = (*BucketArr)[i].InternalIP
		r.ExternalIPArr[i] = (*BucketArr)[i].ExternalIP
		r.InternalPortArr[i] = (*BucketArr)[i].InternalPort
		r.ExternalPortArr[i] = (*BucketArr)[i].ExternalPort
		r.TimestampArr[i] = (*BucketArr)[i].Timestamp
	}
}

func (rpcm *RPCMethod) RPCFindNode(
	cxt context.Context,
	args *ArgsRPCFindNode,
	reply *ReplyRPCFindNode) error {
	var BucketArr []Bucket
	UUID := args.TarUUID
	kb := rpcm.p.kbuckets

	BucketArr = kb.CalClosestKBuckets(UUID)
	reply.FormFromBucketArr(&BucketArr)
	kb.AddBucket(args.FormABucket())

	return nil
}

type ArgsBlockChainSync struct {
	BlockChainLen uint64

	// BlockHead
	VersionArr   []uint8
	IndexArr     []uint64
	TimestampArr []int64
	HashArr      []string
	PreHashArr   []string
	BodyHashArr  []string
	NonceArr     []uint64

	// BlockData
	TxLen        []uint64
	TxHashArr    [][]string
	PreTxHashArr [][][]string
	AddressArr   [][][]string
	BalanceArr   [][][]float64
	SignatureArr [][]string
	NumArr       [][]uint64
}

func NewArgsBlockChainSync(pow *Model.ProofOfWork) *ArgsBlockChainSync {
	var a ArgsBlockChainSync
	var BC *Model.BlockChain
	var i, j, k uint64
	var version uint8

	a = ArgsBlockChainSync{}
	BC = pow.BC

	a.BlockChainLen = BC.Len
	if BC.Len < 1 {
		return &a
	}
	version = Model.Version

	// BlockHead
	a.VersionArr = make([]uint8, a.BlockChainLen, a.BlockChainLen)
	a.IndexArr = make([]uint64, a.BlockChainLen, a.BlockChainLen)
	a.TimestampArr = make([]int64, a.BlockChainLen, a.BlockChainLen)
	a.HashArr = make([]string, a.BlockChainLen, a.BlockChainLen)
	a.PreHashArr = make([]string, a.BlockChainLen, a.BlockChainLen)
	a.BodyHashArr = make([]string, a.BlockChainLen, a.BlockChainLen)
	a.NonceArr = make([]uint64, a.BlockChainLen, a.BlockChainLen)

	// BlockData
	a.TxLen = make([]uint64, a.BlockChainLen, a.BlockChainLen)
	a.TxHashArr = make([][]string, a.BlockChainLen, a.BlockChainLen)
	a.PreTxHashArr = make([][][]string, a.BlockChainLen, a.BlockChainLen)
	a.AddressArr = make([][][]string, a.BlockChainLen, a.BlockChainLen)
	a.BalanceArr = make([][][]float64, a.BlockChainLen, a.BlockChainLen)
	a.SignatureArr = make([][]string, a.BlockChainLen, a.BlockChainLen)
	a.NumArr = make([][]uint64, a.BlockChainLen, a.BlockChainLen)

	for i = 0; i < a.BlockChainLen; i++ {
		// BlockHead
		a.VersionArr[i] = version
		a.IndexArr[i] = BC.HeadArr[i].Index
		a.TimestampArr[i] = BC.HeadArr[i].Timestamp
		a.HashArr[i] = BC.HeadArr[i].Hash
		a.PreHashArr[i] = BC.HeadArr[i].PreHash
		a.BodyHashArr[i] = BC.HeadArr[i].BodyHash
		a.NonceArr[i] = BC.HeadArr[i].Nonce

		// BlockData
		a.TxLen[i] = (*BC.BodyArr)[i].Len
		a.TxHashArr[i] = make([]string, a.TxLen[i], a.TxLen[i])
		a.PreTxHashArr[i] = make([][]string, a.TxLen[i], a.TxLen[i])
		a.AddressArr[i] = make([][]string, a.TxLen[i], a.TxLen[i])
		a.BalanceArr[i] = make([][]float64, a.TxLen[i], a.TxLen[i])
		a.SignatureArr[i] = make([]string, a.TxLen[i], a.TxLen[i])
		a.NumArr[i] = make([]uint64, a.TxLen[i], a.TxLen[i])

		for j = 0; j < a.TxLen[i]; j++ {
			a.TxHashArr[i][j] = (*BC.BodyArr)[i].TxArr[j].Hash
			a.SignatureArr[i][j] = (*BC.BodyArr)[i].TxArr[j].Signature
			a.NumArr[i][j] = (*BC.BodyArr)[i].TxArr[j].Num

			a.PreTxHashArr[i][j] = make([]string, a.NumArr[i][j], a.NumArr[i][j])
			a.AddressArr[i][j] = make([]string, a.NumArr[i][j], a.NumArr[i][j])
			a.BalanceArr[i][j] = make([]float64, a.NumArr[i][j], a.NumArr[i][j])

			for k = 0; k < a.NumArr[i][j]; k++ {
				a.PreTxHashArr[i][j][k] = (*BC.BodyArr)[i].TxArr[j].PreTxHash[k]
				a.AddressArr[i][j][k] = (*BC.BodyArr)[i].TxArr[j].UTXOArr[k].Address
				a.BalanceArr[i][j][k] = (*BC.BodyArr)[i].TxArr[j].UTXOArr[k].Balance
			}
		}
	}

	return &a
}

func (a *ArgsBlockChainSync) FormPOW(pow *Model.ProofOfWork) {
	var BC *Model.BlockChain
	var i, j, k uint64
	var BlockDataArr []Model.BlockData

	BC = pow.BC
	if !pow.VerifyBlockChain() {
		return
	}

	BC.Len = a.BlockChainLen
	BC.HeadArr = make([]Model.BlockHead, Model.MaxPossibleBlockChainLen, Model.MaxPossibleBlockChainLen)
	BlockDataArr = make([]Model.BlockData, Model.MaxPossibleBlockChainLen, Model.MaxPossibleBlockChainLen)
	BC.BodyArr = &BlockDataArr

	for i = 0; i < a.BlockChainLen; i++ {
		// BlockHead
		BC.HeadArr[i].SetVersion(a.VersionArr[i])
		BC.HeadArr[i].Index = a.IndexArr[i]
		BC.HeadArr[i].Timestamp = a.TimestampArr[i]
		BC.HeadArr[i].Hash = a.HashArr[i]
		BC.HeadArr[i].PreHash = a.PreHashArr[i]
		BC.HeadArr[i].BodyHash = a.BodyHashArr[i]
		BC.HeadArr[i].Nonce = a.NonceArr[i]

		// BlockData
		(*BC.BodyArr)[i].Len = a.TxLen[i]
		(*BC.BodyArr)[i].TxArr = make([]Model.Transaction, a.TxLen[i], a.TxLen[i])

		for j = 0; j < a.TxLen[i]; j++ {
			(*BC.BodyArr)[i].TxArr[j].Hash = a.TxHashArr[i][j]
			(*BC.BodyArr)[i].TxArr[j].Signature = a.SignatureArr[i][j]
			(*BC.BodyArr)[i].TxArr[j].Num = a.NumArr[i][j]

			(*BC.BodyArr)[i].TxArr[j].PreTxHash = make([]string, a.NumArr[i][j], a.NumArr[i][j])
			(*BC.BodyArr)[i].TxArr[j].UTXOArr = make([]Model.UTXO, a.NumArr[i][j], a.NumArr[i][j])

			for k = 0; k < a.NumArr[i][j]; k++ {
				(*BC.BodyArr)[i].TxArr[j].PreTxHash[k] = a.PreTxHashArr[i][j][k]
				(*BC.BodyArr)[i].TxArr[j].UTXOArr[k].Address = a.AddressArr[i][j][k]
				(*BC.BodyArr)[i].TxArr[j].UTXOArr[k].Balance = a.BalanceArr[i][j][k]
			}
		}
	}
}

func (rpcm *RPCMethod) RPCBlockChainSync(
	cxt context.Context,
	args *ArgsBlockChainSync,
	reply *ArgsBlockChainSync) error {
	var LocalPOW *Model.ProofOfWork
	var replyPtr *ArgsBlockChainSync
	replyPtr = &ArgsBlockChainSync{}
	replyPtr.BlockChainLen = 0

	LocalPOW = KCI.GetPOW()
	if args.BlockChainLen == LocalPOW.BC.Len {
		return nil
	}

	// return local longer blockchain to remote peer
	if args.BlockChainLen < LocalPOW.BC.Len {
		replyPtr = NewArgsBlockChainSync(LocalPOW)
		*reply = *replyPtr
		return nil
	}

	// Update local blockchain from remote peer
	args.FormPOW(LocalPOW)
	LocalPOW.ClearTxPool()

	return nil
}
