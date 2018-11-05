package KademliaInformationSystem

import (
	"fmt"
	"net"
)

type KBuckets struct {
	BitSpaceNum uint32
	K           uint32
	alpha       uint32
	Buckets     *[][]Bucket
	BucketsCnt  []uint32
	peer        *Peer

	u Utility
}

type Bucket struct {
	UUID         string
	InternalIP   net.IP
	ExternalIP   net.IP
	InternalPort uint32
	ExternalPort uint32
	Timestamp    int64
}

func (b *Bucket) FormInternalAddr() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d",
		b.InternalIP[0], b.InternalIP[1],
		b.InternalIP[2], b.InternalIP[3],
		b.InternalPort)
}

func (b *Bucket) FormExternalAddr() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d",
		b.ExternalIP[0], b.ExternalIP[1],
		b.ExternalIP[2], b.ExternalIP[3],
		b.ExternalPort)
}

func NewKBuckets() *KBuckets {
	kb := KBuckets{}
	kb.u = Utility{}

	kb.BitSpaceNum = 160
	kb.K = 20
	kb.alpha = 3
	Buckets := make([]([]Bucket), kb.BitSpaceNum, kb.BitSpaceNum)
	for i, _ := range Buckets {
		Buckets[i] = make([]Bucket, kb.K, kb.K)
	}
	kb.Buckets = &Buckets
	kb.BucketsCnt = make([]uint32, kb.BitSpaceNum, kb.BitSpaceNum)

	return &kb
}

func (kb *KBuckets) CalClosestKBuckets(UUID string) []Bucket {
	var BucketArr []Bucket
	var i, j uint32

	BitIdx := kb.u.CalBitIdxStr(kb.peer.UUID, UUID)
	TarBucketArr := (*kb.Buckets)[BitIdx]
	LocBucketArr := TarBucketArr[0:kb.BucketsCnt[BitIdx]]

	if kb.BucketsCnt[BitIdx] == kb.K {
		BucketArr = LocBucketArr
		return BucketArr
	}

	BucketArrCnt := kb.BucketsCnt[BitIdx]
	BucketArr = make([]Bucket, kb.K, kb.K)
	for i = 0; i < BucketArrCnt; i++ {
		BucketArr[i] = LocBucketArr[i]
	}

	for i = BitIdx + 1; i < kb.BitSpaceNum; i++ {
		if BucketArrCnt+kb.BucketsCnt[i] >= kb.K {
			// TODO:not exactly the closest k-buckets yet
			for j = BucketArrCnt; j < kb.K; j++ {
				BucketArr[j] = (*kb.Buckets)[i][j]
			}
			BucketArrCnt = kb.K
			break
		}

		for j = BucketArrCnt; j < BucketArrCnt+kb.BucketsCnt[i]; j++ {
			BucketArr[j] = (*kb.Buckets)[i][j-BucketArrCnt]
		}
		BucketArrCnt += kb.BucketsCnt[i]
	}
	if BucketArrCnt >= kb.K {
		return BucketArr
	}

	for i = BitIdx - 1; i < BitIdx && i >= 0; i-- {
		if BucketArrCnt+kb.BucketsCnt[i] >= kb.K {
			// TODO:not exactly the closest k-buckets yet
			for j = BucketArrCnt; j < kb.K; j++ {
				BucketArr[j] = (*kb.Buckets)[i][j]
			}
			break
		}

		for j = BucketArrCnt; j < BucketArrCnt+kb.BucketsCnt[i]; j++ {
			BucketArr[j] = (*kb.Buckets)[i][j-BucketArrCnt]
		}
		BucketArrCnt += kb.BucketsCnt[i]
	}

	return BucketArr[0:BucketArrCnt]
}

func (kb *KBuckets) AddBucket(bucket Bucket) {
	var i int

	BitIdx := kb.u.CalBitIdxStr(kb.peer.UUID, bucket.UUID)

	for i = 0; i < int(kb.BucketsCnt[BitIdx]); i++ {
		if (*kb.Buckets)[BitIdx][i].UUID == bucket.UUID {
			return
		}
	}

	if kb.BucketsCnt[BitIdx] < kb.K {
		(*kb.Buckets)[BitIdx][kb.BucketsCnt[BitIdx]] = bucket
		kb.BucketsCnt[BitIdx]++

		// if Kci != nil {
		// 	Kci.ShowKBuckets(kb)
		// }
		return
	}

	EvictedIdx := -1
	peer := kb.peer
	for i = 0; i < int(kb.K); i++ {
		if Kad.Scenario == RealScenario {
			if peer.PingAddr(
				(*kb.Buckets)[BitIdx][i].ExternalIP,
				(*kb.Buckets)[BitIdx][i].ExternalPort) {
				continue
			}
		} else {
			if Kad.Ping((*kb.Buckets)[BitIdx][i].UUID) {
				continue
			}
		}
		EvictedIdx = i
		break
	}
	if EvictedIdx < 0 {
		return
	}

	for i = EvictedIdx; i < int(kb.K-1); i++ {
		(*kb.Buckets)[BitIdx][i] =
			(*kb.Buckets)[BitIdx][i+1]
	}
	(*kb.Buckets)[BitIdx][kb.K-1] = bucket
}

func (kb *KBuckets) NodeLookUp(
	UUID string) (
	FoundBucket bool,
	bucket Bucket,
	AlphaSearchCnt uint32) {
	var i, j, BucketArrCnt, BucketPoolCnt int
	var MinDis, CurMinDis float64
	var LocBucketArr, BucketArr, BucketPool, FindResult []Bucket
	var u Utility
	var p *Peer

	u = Utility{}
	p = (*kb).peer
	FoundBucket = false
	bucket = Bucket{}
	AlphaSearchCnt = 0
	LocBucketArr = kb.CalClosestKBuckets(UUID)

	BucketArr = make([]Bucket, kb.K, kb.K)
	BucketPool = make([]Bucket, kb.K*kb.K, kb.K*kb.K)

	BucketArrCnt = len(LocBucketArr)
	for i = 0; i < BucketArrCnt; i++ {
		BucketArr[i] = LocBucketArr[i]
	}

	MinDis = 1.797e+308 // realmax

	for {
		AlphaSearchCnt++

		BucketPoolCnt = 0
		for i = 0; i < BucketArrCnt; i++ {
			if Kad.Scenario == RealScenario {
				taddr := BucketArr[i].FormExternalAddr()
				fmt.Println(taddr)

				if kb.peer.CallRPCPing(BucketArr[i].FormExternalAddr()) {
					_, FindResult = p.CallRPCFindNode(
						BucketArr[i].FormExternalAddr(), UUID)
				} else {
					continue
				}
			} else {
				if Kad.Ping(BucketArr[i].UUID) {
					FindResult = Kad.FindNode(BucketArr[i].UUID, UUID, *(p.FormABucket()))
				} else {
					continue
				}
			}

			for j = BucketPoolCnt; j < BucketPoolCnt+len(FindResult); j++ {
				BucketPool[j] = FindResult[j-BucketPoolCnt]
			}
			BucketPoolCnt += len(FindResult)
		}
		if BucketPoolCnt == 0 {
			break
		}
		for i = 0; i < BucketPoolCnt; i++ {
			kb.AddBucket(BucketPool[i])
		}

		SortIdx := u.SortBucketArr(BucketPool[0:BucketPoolCnt], UUID)
		if BucketPoolCnt > int((*kb).K) {
			SortIdx = SortIdx[0:kb.K]
		}
		BucketArrCnt = len(SortIdx)
		BucketArr = make([]Bucket, BucketArrCnt, BucketArrCnt)
		for i = 0; i < BucketArrCnt; i++ {
			BucketArr[i] = BucketPool[SortIdx[i]]
		}

		CurMinDis = u.CalXorDisStr(BucketPool[SortIdx[0]].UUID, UUID)
		bucket = BucketPool[SortIdx[0]]
		if CurMinDis == 0 {
			FoundBucket = true
			break
		}

		if MinDis <= CurMinDis {
			break
		}

		MinDis = CurMinDis
	}

	return FoundBucket, bucket, AlphaSearchCnt
}

func (kb *KBuckets) CalBucketNum() uint16 {
	var BucketNum, i uint32

	for i = 0; i < kb.BitSpaceNum; i++ {
		BucketNum += kb.BucketsCnt[i]
	}

	return uint16(BucketNum)
}

func (kb *KBuckets) GetAllAddr() []string {
	var AllAddr []string
	var i, j, AddrNumCnt uint32
	AllAddr = make([]string, kb.BitSpaceNum*kb.K, kb.BitSpaceNum*kb.K)

	AddrNumCnt = 0
	for i = 0; i < kb.BitSpaceNum; i++ {
		for j = 0; j < kb.BucketsCnt[i]; j++ {
			AllAddr[AddrNumCnt] = xu.IP2Str(
				(*kb.Buckets)[i][j].ExternalIP,
				(*kb.Buckets)[i][j].ExternalPort)
			AddrNumCnt++
		}
	}
	AllAddr = AllAddr[0:AddrNumCnt]

	return AllAddr
}
