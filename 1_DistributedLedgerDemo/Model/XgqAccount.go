package Model

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const (
	EP521 = iota
	EP384
	EP256
	EP224
)

type XgqAccount struct {
	PublicKeyHash string
	Nickname      string

	privateKey ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey

	Balance  float64
	UtxoHash string

	SealedCryptoObj SealedCrypto
}

// TO be delete
func (XgqA *XgqAccount) GetPrivateKey() (privateKey ecdsa.PrivateKey) {
	return XgqA.privateKey
}

// TO be delete
func (XgqA *XgqAccount) SetPrivateKey(privateKey ecdsa.PrivateKey) {
	XgqA.privateKey = privateKey
}

func (XgqA *XgqAccount) GenRandAccount() {
	var RdInt64V uint64

	// neither "crypto/rand" nor "math/rand" is not truely random
	// need a true random seed

	// RdInt64V = rand.Uint64()
	Seed := make([]byte, 20)
	RdInt, _ := rand.Read(Seed)
	RdInt64V = uint64(RdInt)

	XgqA.GenAccount(RdInt64V)
}

func (XgqA *XgqAccount) GenAccount(RdNum uint64) (Addr string, X, Y, D big.Int, CurType int) {
	h := sha256.New()

	h.Write([]byte(strconv.FormatUint(RdNum, 10)))
	hashed1 := h.Sum(nil)

	h.Write([]byte(hashed1))
	hashed2 := h.Sum(nil)

	XgqA.GenKeyPair(hex.EncodeToString(hashed1) + hex.EncodeToString(hashed2))

	h = sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", XgqA.PublicKey.X)))
	hashed3 := h.Sum(nil)

	XgqA.PublicKeyHash = hex.EncodeToString(hashed3)

	XgqA.Balance = 0
	XgqA.UtxoHash = ""

	XgqA.SealedCryptoObj = SealedCrypto{}

	return XgqA.PublicKeyHash, *XgqA.privateKey.X,
		*XgqA.privateKey.Y, *XgqA.privateKey.D, EP521
}

func (XgqA *XgqAccount) GenKeyPair(RdStr string) {
	var err error
	var CurveType int
	var curve elliptic.Curve

	CurveType = EP521

	switch CurveType {
	case EP521:
		curve = elliptic.P521()
	case EP384:
		curve = elliptic.P384()
	case EP256:
		curve = elliptic.P256()
	default:
		curve = elliptic.P224()
	}

	privateKey, err := ecdsa.GenerateKey(curve, strings.NewReader(RdStr))
	XgqA.privateKey = *privateKey
	if err != nil {
		return
	}

	XgqA.PublicKey = XgqA.privateKey.PublicKey
}

func (XgqA *XgqAccount) Sign(text string) (string, error) {
	return XgqA.SealedCryptoObj.Sign(text, &XgqA.privateKey)
}

func (XgqA *XgqAccount) Verify(text, passwd string) (bool, error) {
	return XgqA.SealedCryptoObj.Verify(text, passwd, &XgqA.PublicKey)
}

func (XgqA *XgqAccount) Pay(
	AdrrArr []string,
	LastValueArr []float64,
	PayValueArr []float64,
	UtxoHashArr []string) (
	Tx Transaction, err error) {

	err = nil

	Len := len(AdrrArr)
	Tx.PreTxHash = make([]string, Len, Len)
	Tx.UTXOArr = make([]UTXO, Len, Len)

	Tx.PreTxHash[0] = XgqA.UtxoHash
	Tx.UTXOArr[0].Address = XgqA.PublicKeyHash

	for i := 1; i < Len; i++ {
		Tx.PreTxHash[i] = UtxoHashArr[i]
		Tx.UTXOArr[i].Address = AdrrArr[i]
		Tx.UTXOArr[i].Balance = LastValueArr[i] + PayValueArr[i]

		XgqA.Balance -= PayValueArr[i]
	}

	if XgqA.Balance < 0 {
		err = errors.New("Not Balanced!")
		return Tx, err
	}

	Tx.UTXOArr[0].Balance = XgqA.Balance
	Tx.Num = uint64(Len)

	Tx.CalHash()
	Tx.Signature, _ = XgqA.Sign(Tx.Hash)

	return Tx, err
}
