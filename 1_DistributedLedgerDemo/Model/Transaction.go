package Model

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"hash"
)

type Transaction struct {
	Hash string

	// First address pays

	PreTxHash []string
	UTXOArr   []UTXO
	Signature string
	Num       uint64
}

type UTXO struct {
	Address string
	Balance float64
}

func (Tx *Transaction) CalHash() {
	h := sha256.New()

	Tx.writeBytes(h)

	Tx.Hash = hex.EncodeToString(h.Sum(nil))
}

func (Tx *Transaction) writeBytes(h hash.Hash) {
	var i uint64

	for i = 0; i < Tx.Num; i++ {
		h.Write([]byte(Tx.PreTxHash[i]))

		Tx.UTXOArr[i].writeBytes(h)
	}
	h.Write([]byte(Tx.Signature))

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, Tx.Num)
	h.Write(buffer.Bytes())
}

func (u *UTXO) writeBytes(h hash.Hash) {
	h.Write([]byte(u.Address))

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, u.Balance)
	h.Write(buffer.Bytes())
}
