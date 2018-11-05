package Model

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"math/big"
	"time"
)

type Block struct {
	Head BlockHead
	Body *BlockData
}

type BlockHead struct {
	version   uint8
	Index     uint64
	Timestamp int64
	Hash      string
	PreHash   string
	BodyHash  string
	Nonce     uint64
}

type BlockData struct {
	TxArr []Transaction
	Len   uint64
}

func (b *Block) GetGenesisBlock() (gb Block) {

	Head := &gb.Head
	gb.Body = &BlockData{}

	Head.version = Version
	Head.Index = 1

	gbTime, _ := time.Parse("2006-01-02 15:04:05", "2018-07-16 11:12:13")
	Head.Timestamp = gbTime.Unix()

	Head.PreHash = ""
	Head.BodyHash = gb.Body.CalHash()
	Head.Nonce = 0

	h := sha256.New()
	Head.writeBytes(h)
	h.Write([]byte("This is genesis block of the XGQ blockchain demo."))
	Head.Hash = hex.EncodeToString(h.Sum(nil))

	return gb
}

func (b *Block) CalHash() {
	h := sha256.New()

	b.Head.BodyHash = b.Body.CalHash()
	b.Head.writeBytes(h)

	b.Head.Hash = hex.EncodeToString(h.Sum(nil))
}

func (b *Block) CalHash4Mine() big.Int {
	h := sha256.New()
	b.Head.writeBytes(h)
	bytes := h.Sum(nil)

	var hashInt big.Int
	hashInt.SetBytes(bytes[:])

	return hashInt
}

func (h *BlockHead) SetVersion(v uint8) {
	h.version = v
}

func (h *BlockHead) Initilization(
	Idx uint64,
	BodyHash string,
	PreHash string) {
	h.version = Version
	h.Index = Idx

	h.Timestamp = time.Now().Unix()

	h.BodyHash = BodyHash
	h.PreHash = PreHash
	h.Nonce = 0
}

func (Head *BlockHead) writeBytes(h hash.Hash) {
	h.Write([]byte{Head.version})

	ByteArr := make([]byte, 8)
	binary.BigEndian.PutUint64(ByteArr, Head.Index)
	h.Write(ByteArr)

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, Head.Timestamp)
	binary.Write(&buffer, binary.BigEndian, Head.Nonce)
	h.Write(buffer.Bytes())

	h.Write([]byte(Head.PreHash))
	h.Write([]byte(Head.BodyHash))
}

func (bd *BlockData) Initilization(TxArr []Transaction, Len uint64) {
	// bd.TxArr = TxArr
	bd.Len = Len

	bd.TxArr = make([]Transaction, len(TxArr), len(TxArr))

	// go func() {
	// 	for i := 0; i < len(TxArr); i++ {
	// 		bd.TxArr[i] = TxArr[i]
	// 	}
	// }()

	for i := 0; i < len(TxArr); i++ {
		bd.TxArr[i] = TxArr[i]
	}
}

func (Body *BlockData) writeBytes(h hash.Hash) {
	TxArr := Body.TxArr
	h.Write([]byte(fmt.Sprintf("%d", Body.Len)))
	var i uint64
	for i = 0; i < Body.Len; i++ {
		TxArr[i].writeBytes(h)
	}
}

func (Body *BlockData) CalHash() string {
	h := sha256.New()
	Body.writeBytes(h)
	return hex.EncodeToString(h.Sum(nil))
}

func (bd *BlockData) AddTrasaction(Tx Transaction) {
	bd.TxArr[bd.Len] = Tx
	bd.Len++
}

func (bd *BlockData) QueryBlance(Addr *string) float64 {
	var i uint64
	var UTXOArr []UTXO

	if bd.Len < 1 {
		return -1
	}

	for i = bd.Len - 1; i < bd.Len; i-- {
		UTXOArr = bd.TxArr[i].UTXOArr

		for j := 0; j < len(UTXOArr); j++ {
			if UTXOArr[j].Address == *Addr {
				return UTXOArr[j].Balance
			}
		}
	}

	return -1
}
