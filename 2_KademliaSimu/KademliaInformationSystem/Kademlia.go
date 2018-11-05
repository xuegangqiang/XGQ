package KademliaInformationSystem

import (
	"XGQ/1_DistributedLedgerDemo/Model"
	"XGQ/XgqUtility"
	"fmt"
	"math/rand"
)

var Kad *Kademlia

const (
	SimulateScenario = iota
	RealScenario
)

var xu XgqUtility.Utility

var KCI *KisCtrlInterface

type KisCtrlInterface struct {
	Report2Ctrl func(msg string)
	GetPOW      func() *Model.ProofOfWork
}

type Kademlia struct {
	PeerArr  *[]Peer
	PeerCnt  uint64
	Scenario int

	UUIDArr []string
}

func NewKademlia() *Kademlia {
	var MaxPeerNum uint64

	k := Kademlia{}

	MaxPeerNum = 999
	PeerArr := make([]Peer, MaxPeerNum, MaxPeerNum)
	k.PeerArr = &PeerArr
	k.PeerCnt = 0

	k.Scenario = RealScenario

	return &k
}

func (k *Kademlia) AddAPeer() {
	NewPeer := (*NewPeer())

	(*k.PeerArr)[k.PeerCnt] = NewPeer
	k.PeerCnt++

	if k.PeerCnt == 1 {
		return
	}

	RefNodeNum := 1
	RefNodeIdx := make([]int, RefNodeNum, RefNodeNum)
	RefNodeArr := make([]Peer, RefNodeNum, RefNodeNum)
	RefNodeArrPtr := &RefNodeArr

	for i, _ := range RefNodeIdx {
		RefNodeIdx[i] = rand.Intn(int(k.PeerCnt) - 1)
		RefNodeArr[i] = (*k.PeerArr)[RefNodeIdx[i]]
	}

	NewPeer.JoinNetWork(RefNodeArrPtr)
	(*k.PeerArr)[k.PeerCnt-1] = NewPeer
}

func (k *Kademlia) DoSimu() {
	var InitPeerNum, i uint64

	k.Scenario = SimulateScenario

	InitPeerNum = 990
	k.UUIDArr = make([]string, InitPeerNum, InitPeerNum)
	for i = 0; i < InitPeerNum; i++ {
		k.AddAPeer()

		k.UUIDArr[i] = (*k.PeerArr)[i].UUID
	}
}

func (k *Kademlia) Ping(UUID string) bool {
	for _, p := range *(k.PeerArr) {
		if p.UUID == UUID {
			return true
		}
	}

	return false
}

func (k *Kademlia) FindNode(
	DstPeerUUID, TarUUID string,
	SourPeerBucket Bucket) []Bucket {
	var i uint64
	var BucketArr []Bucket
	var DstKB KBuckets

	for i = 0; i < k.PeerCnt; i++ {
		if (*k.PeerArr)[i].UUID != DstPeerUUID {
			continue
		}
		DstKB = (*k.PeerArr)[i].kbuckets

		BucketArr = DstKB.CalClosestKBuckets(TarUUID)
		DstKB.AddBucket(SourPeerBucket)

		break
	}

	return BucketArr
}

func (k *Kademlia) VerifySearchAbility() {
	var i, FoundCnt uint64
	var AlphaSearchCnt uint32
	var RandIdx int
	var p Peer
	var kb KBuckets
	var FoundBucket bool
	var info string

	k.PeerOfflineSimu()

	PeerCnt := k.PeerCnt
	PeerArr := (*k.PeerArr)[0:PeerCnt]

	// RdIdxArr=randi(PeerCnt,[PeerCnt,1]);
	FoundCnt = 0
	for i = 0; i < PeerCnt; i++ {
		p = PeerArr[i]
		kb = p.kbuckets

		RandIdx = rand.Intn(int(PeerCnt) - 1)

		// [FoundBucket,~,AlphaSearchCnt]=...
		// 	kb.NodeLookUp(PeerArr(RdIdxArr(i)).UUID);

		// aa := len(PeerArr)
		// aa1 := PeerArr[RandIdx].UUID
		// if aa == 0 && aa1 == "" {
		// }

		FoundBucket, _, AlphaSearchCnt = kb.NodeLookUp(
			PeerArr[RandIdx].UUID)

		info = fmt.Sprintf("%2d,AlphaSearchCnt:%3d,result:", i, AlphaSearchCnt)
		if FoundBucket {
			FoundCnt = FoundCnt + 1
			info = fmt.Sprintf("%strue ,", info)
		} else {
			info = fmt.Sprintf("%sfalse,", info)
		}
		// fprintf('BucketsNum:%2d\n',sum(p.kbuckets.BucketsCnt));
		info = fmt.Sprintf("%sBucketsNum:%2d", info, kb.CalBucketNum())
		fmt.Println(info)
	}
	info = fmt.Sprintf("%3d/%3d search failed", PeerCnt-FoundCnt, PeerCnt)
	fmt.Println(info)
}

func (k *Kademlia) PeerOfflineSimu() {
	var PeerCnt0, i, NewPeerCnt uint64
	var OfflinePeerRatio float64
	var OfflineIdx []bool
	var NewPeerArr []Peer

	OfflinePeerRatio = 0.5

	PeerCnt0 = k.PeerCnt
	OfflineIdx = make([]bool, PeerCnt0, PeerCnt0)

	for i = 0; i < PeerCnt0; i++ {

		if float64(rand.Intn(10000))/10000 > OfflinePeerRatio {
			OfflineIdx[i] = true
			continue
		}

		k.PeerCnt--
	}

	NewPeerArr = make([]Peer, k.PeerCnt, k.PeerCnt)
	NewPeerCnt = 0
	for i = 0; i < PeerCnt0; i++ {
		if !OfflineIdx[i] {
			continue
		}

		NewPeerArr[NewPeerCnt] = (*k.PeerArr)[i]
		NewPeerCnt++
	}

	k.PeerArr = &NewPeerArr
}
