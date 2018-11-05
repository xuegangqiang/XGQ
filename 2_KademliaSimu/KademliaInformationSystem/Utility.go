package KademliaInformationSystem

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type Utility struct {
}

func (u *Utility) CalXorDisStr(Str1, Str2 string) float64 {
	XorDis := u.CalXorDis([]byte(Str1), []byte(Str2))
	return XorDis
}

func (u *Utility) CalXorDis(Arr1, Arr2 []uint8) float64 {
	var XorDis float64

	Dis := u.ByteArrXor(Arr1, Arr2)

	ItemNum := len(Dis)
	XorDis = 0
	for i := 0; i < ItemNum; i++ {
		XorDis += float64(256^i) * float64(Dis[i])
	}

	return XorDis
}

func (u *Utility) CalBitIdxStr(UUIDStr, TarIDStr string) uint32 {
	var UUID, TarID []uint8
	var i int

	UUIDStr = strings.Replace(UUIDStr, "-", "", -1)
	TarIDStr = strings.Replace(TarIDStr, "-", "", -1)

	UUID, _ = hex.DecodeString(UUIDStr)
	TarID, _ = hex.DecodeString(TarIDStr)
	// UUID = []byte(UUIDStr)
	// TarID = []byte(TarIDStr)

	L1 := len(UUID)
	L2 := len(TarID)
	if L1 > L2 {
		TarIDExt := make([]uint8, L1, L1)
		for i = L1 - L2; i < L1; i++ {
			TarIDExt[i] = TarID[i-(L1-L2)]
		}
		TarID = TarIDExt
	}

	if L1 < L2 {
		UUIDExt := make([]uint8, L2, L2)
		for i = L2 - L1; i < L2; i++ {
			UUIDExt[i] = UUID[i-(L2-L1)]
		}
		UUID = UUIDExt
	}

	BitIdx := u.CalBitIdx(UUID, TarID)

	return BitIdx
}

func (u *Utility) CalBitIdx(UUID, TarID []uint8) uint32 {
	var BitIdx uint32
	var BinStr string

	ByteNum := len(UUID)
	BitNum := ByteNum * 8
	Dis := u.ByteArrXor(UUID, TarID)

	BitIdx = uint32(BitNum)
	for i := 0; i < ByteNum; i++ {
		BinStr = fmt.Sprintf("%08b", Dis[i])

		for j := 0; j < 8; j++ {
			if BinStr[j] == '0' {
				continue
			}

			BitIdx = uint32(i*8 + j)
			return BitIdx
		}
	}

	return BitIdx
}

func (u *Utility) ByteArrXor(Arr1, Arr2 []uint8) []uint8 {
	ByteNum := len(Arr1)
	Dis := make([]uint8, ByteNum, ByteNum)

	for i, _ := range Dis {
		Dis[i] = Arr1[i] ^ Arr2[i]
	}

	return Dis
}

func (u *Utility) SortBucketArr(BucketArr []Bucket, ID string) []uint32 {
	ItemNum := len(BucketArr)
	XorDisArr := make([]float64, ItemNum, ItemNum)

	for i := 0; i < ItemNum; i++ {
		XorDisArr[i] = u.CalXorDisStr(BucketArr[i].UUID, ID)
	}

	SortIdx := u.SortWithIdx(XorDisArr)
	return SortIdx
}

func (u *Utility) SortWithIdx(Arr []float64) []uint32 {
	ItemNum := len(Arr)
	m := make(map[float64]uint32, ItemNum)
	for i, v := range Arr {
		m[v] = uint32(i)
	}

	sort.Float64s(Arr)
	Idx := make([]uint32, ItemNum, ItemNum)
	for i, v := range Arr {
		Idx[i] = m[v]
	}

	return Idx
}
