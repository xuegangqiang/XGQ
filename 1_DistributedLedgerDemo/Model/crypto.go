package Model

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
)

const (
	ivDefValue = "0102030405060708"
)

type SealedCrypto struct {
}

// ecdsa stuff
func (sc *SealedCrypto) Sign(text string, PrivateKey *ecdsa.PrivateKey) (string, error) {
	var RdInt64V uint64
	RdInt64V = rand.Uint64()
	h := sha256.New()
	h.Write([]byte(strconv.FormatUint(RdInt64V, 10)))
	hashed := h.Sum(nil)
	RdStr := hex.EncodeToString(hashed)

	r, s, err := ecdsa.Sign(strings.NewReader(RdStr), PrivateKey, []byte(text))

	if err != nil {
		return "", err
	}
	rt, err := r.MarshalText()
	if err != nil {
		return "", err
	}
	st, err := s.MarshalText()
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	_, err = w.Write([]byte(string(rt) + "+" + string(st)))
	if err != nil {
		return "", err
	}
	w.Flush()
	return hex.EncodeToString(b.Bytes()), nil
}

func (sc *SealedCrypto) getSign(
	text, byterun []byte) (rint, sint big.Int, err error) {
	r, err := gzip.NewReader(bytes.NewBuffer(byterun))
	if err != nil {
		err = errors.New("decode error," + err.Error())
		return
	}
	defer r.Close()
	buf := make([]byte, 1024)
	count, err := r.Read(buf)
	if err != nil {
		fmt.Println("decode =", err)
		err = errors.New("decode read error," + err.Error())
		return
	}
	rs := strings.Split(string(buf[:count]), "+")
	if len(rs) != 2 {
		err = errors.New("decode fail")
		return
	}
	err = rint.UnmarshalText([]byte(rs[0]))
	if err != nil {
		err = errors.New("decrypt rint fail, " + err.Error())
		return
	}
	err = sint.UnmarshalText([]byte(rs[1]))
	if err != nil {
		err = errors.New("decrypt sint fail, " + err.Error())
		return
	}
	return
}

func (sc *SealedCrypto) Verify(
	text, passwd string, PublicKey *ecdsa.PublicKey) (bool, error) {
	byterun, err := hex.DecodeString(passwd)
	if err != nil {
		return false, err
	}
	rint, sint, err := sc.getSign([]byte(text), byterun)
	if err != nil {
		return false, err
	}
	result := ecdsa.Verify(PublicKey, []byte(text), &rint, &sint)
	return result, nil
}

// AES stuff
func (sc *SealedCrypto) AesEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("invalid decrypt key")
	}
	blockSize := block.BlockSize()
	plaintext = sc.PKCS5Padding(plaintext, blockSize)
	iv := []byte(ivDefValue)
	blockMode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

func (sc *SealedCrypto) AesDecrypt(ciphertext []byte, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("invalid decrypt key")
	}

	blockSize := block.BlockSize()

	if len(ciphertext) < blockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := []byte(ivDefValue)
	if len(ciphertext)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	blockModel := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plaintext, ciphertext)
	plaintext = sc.PKCS5UnPadding(plaintext)

	return plaintext, nil
}

func (sc *SealedCrypto) PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func (sc *SealedCrypto) PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
