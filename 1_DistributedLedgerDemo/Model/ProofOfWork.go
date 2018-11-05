package Model

import (
	"fmt"
	"math/big"
	"time"
)

type ProofOfWork struct {
	BC *BlockChain

	TargetTime4NewBlock uint16 // Seconds

	Difficulty       uint16 // number of zero bits
	DifficultyArr    []uint16
	DifficultyAdjLim float64

	TxPool    []Transaction
	TxPoolLen uint64
}

func (pow *ProofOfWork) Initilization() {
	pow.BC = &BlockChain{}
	pow.BC.Initilization()

	pow.TargetTime4NewBlock = 10

	pow.Difficulty = 20
	pow.DifficultyArr = make([]uint16, 999, 999)
	pow.DifficultyAdjLim = 1

	pow.ClearTxPool()
}

func (pow *ProofOfWork) Mine() Block {
	bc := pow.BC

	b := Block{}
	b.Body = &BlockData{}

	b.Body.Initilization(pow.TxPool, pow.TxPoolLen)
	b.Head.Initilization(bc.Len+1, b.Body.CalHash(), bc.HeadArr[bc.Len-1].Hash)
	var NonceI uint64
	var hashInt big.Int

	Difficulty := big.NewInt(1)
	Difficulty.Lsh(Difficulty, uint(256-pow.Difficulty))
	fmt.Println("Current difficulty:", pow.Difficulty)

	// TimeS := bc.HeadArr[bc.Len-1].Timestamp // real scenario
	TimeS := time.Now().Unix() // test scenario
	if bc.Len < 2 {
		TimeS = time.Now().Unix()
	}

	NonceI = 0
	for {
		b.Head.Nonce = NonceI
		hashInt = b.CalHash4Mine()
		NonceI++

		if hashInt.Cmp(Difficulty) == -1 {
			break
		}
	}
	b.CalHash()
	TimeE := time.Now().Unix()

	bc.AddBlock(b)
	pow.DifficultyArr[bc.Len-1] = pow.Difficulty
	pow.ClearTxPool()

	DifficultyAdjRatio := float64(pow.TargetTime4NewBlock) / (float64(TimeE) - float64(TimeS))
	if TimeS == TimeE || DifficultyAdjRatio > pow.DifficultyAdjLim {
		DifficultyAdjRatio = pow.DifficultyAdjLim
	}

	// if DifficultyAdjRatio < 1/pow.DifficultyAdjLim {
	// 	DifficultyAdjRatio = 1 / pow.DifficultyAdjLim
	// }

	// pow.Difficulty = uint16(float64(pow.Difficulty) * DifficultyAdjRatio)

	if DifficultyAdjRatio < 1 {
		pow.Difficulty -= 1
	} else {
		pow.Difficulty += uint16(DifficultyAdjRatio)
	}

	return b
}

func (pow *ProofOfWork) Add2TrasactionPool(Tx Transaction, AccoundArr *map[string]XgqAccount) {
	Tx.CalHash()
	pow.TxPool[pow.TxPoolLen] = Tx
	pow.TxPoolLen++

	// TODO:Verify Transaction

	var XgqA XgqAccount
	for _, utxo := range Tx.UTXOArr {
		XgqA = (*AccoundArr)[utxo.Address]
		XgqA.UtxoHash = Tx.Hash
		(*AccoundArr)[utxo.Address] = XgqA
	}
}

func (pow *ProofOfWork) ClearTxPool() {
	pow.TxPool = make([]Transaction, 999, 999)
	pow.TxPoolLen = 0
}

func (pow *ProofOfWork) VerifyBlockChain() bool {
	// TODO
	return true
}
