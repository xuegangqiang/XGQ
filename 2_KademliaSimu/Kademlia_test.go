package main

import (
	kis "XGQ/2_KademliaSimu/KademliaInformationSystem"
	"fmt"
	"testing"
	"time"
)

func TestKademlia(t *testing.T) {
	TimeS := time.Now().UnixNano()

	kis.Kad = kis.NewKademlia()
	Kad := kis.Kad

	Kad.DoSimu()
	Kad.VerifySearchAbility()

	TimeE := time.Now().UnixNano()
	TimeElapsed := float64(TimeE-TimeS) / 1e9
	info := fmt.Sprintf(
		"Time elapsed is %5.2f s,%5.2f s/KPeers",
		TimeElapsed, 1000*TimeElapsed/float64(Kad.PeerCnt))
	fmt.Println(info)
}

func TestRand(t *testing.T) {
}

func TestUtility(t *testing.T) {
}

func TestCalBitIdxStr(t *testing.T) {
	var UUIDStr, TarIDStr string
	var UUID, TarID []uint8

	u := kis.Utility{}

	UUIDStr = "abcfff"
	TarIDStr = "abcdef"

	UUID = []byte(UUIDStr)
	TarID = []byte(TarIDStr)

	BitIdx := u.CalBitIdxStr(UUIDStr, TarIDStr)
	BitIdx = u.CalBitIdx(UUID, TarID)

	info := fmt.Sprintf("%d", BitIdx)
	fmt.Println(info)
}
func TestCalBitIdx(t *testing.T) {
	var UUID, TarID []uint8
	var info string
	UUID = make([]uint8, 2, 2)
	TarID = make([]uint8, 2, 2)

	UUID[0] = 33
	UUID[1] = 4

	TarID[0] = 33
	TarID[1] = 44

	u := kis.Utility{}

	BitIdx := u.CalBitIdx(UUID, TarID)

	info = ""
	for i := 0; i < len(UUID); i++ {
		info += fmt.Sprintf("%08b-", UUID[i])
	}
	fmt.Println(string([]byte(info)[:len(info)-1]))

	info = ""
	for i := 0; i < len(TarID); i++ {
		info += fmt.Sprintf("%08b-", TarID[i])
	}
	fmt.Println(string([]byte(info)[:len(info)-1]))

	info = fmt.Sprintf("%d", BitIdx)
	fmt.Println(info)
}
func TestXOR(t *testing.T) {
	var info string
	var ByteArr1, ByteArr2, ByteArr3 []uint8

	ByteArr1 = make([]uint8, 8, 8)
	ByteArr2 = make([]uint8, 8, 8)
	ByteArr3 = make([]uint8, 8, 8)

	ByteArr1[0] = 5

	ByteArr2[0] = 3

	ByteArr3[0] = ByteArr1[0] ^ ByteArr2[0]
	// ByteArr3 = ByteArr1 ^ ByteArr2

	info = fmt.Sprintf("%08b", ByteArr1[0])

	// info = fmt.Sprintf("%d,%d,%d", ByteArr1[0], ByteArr2[0], ByteArr2[0])
	fmt.Println(info)
}

func TestSortIdx(t *testing.T) {
	var Arr []float64
	Arr = make([]float64, 4, 4)
	u := kis.Utility{}
	Arr[0] = 4
	Arr[1] = 2
	Arr[2] = 1
	Arr[3] = 3

	Idx := u.SortWithIdx(Arr)

	info := fmt.Sprintf("%d,%d,%d,%d", Idx[0], Idx[1], Idx[2], Idx[3])
	fmt.Println(info)
}
