package Model

const Version uint8 = 1
const MaxPossibleBlockChainLen uint64 = 999

type BlockChain struct {
	HeadArr []BlockHead
	BodyArr *[]BlockData

	Len uint64
}

func (bc *BlockChain) Initilization() {
	bc.HeadArr = make([]BlockHead, MaxPossibleBlockChainLen, MaxPossibleBlockChainLen)
	BodyArr := make([]BlockData, MaxPossibleBlockChainLen, MaxPossibleBlockChainLen)
	bc.BodyArr = &BodyArr

	bc.Len = 1

	b := Block{}
	b = b.GetGenesisBlock()
	bc.HeadArr[0] = b.Head
	BodyArr[0] = BlockData{}
}

func (bc *BlockChain) AddBlock(b Block) {
	bc.HeadArr[bc.Len] = b.Head

	BodyArr := *bc.BodyArr

	BodyArr[bc.Len] = *b.Body

	bc.BodyArr = &BodyArr

	bc.Len++
}

func (bc *BlockChain) QueryBlance(Addr string) float64 {
	var i uint64
	var Blance, cb float64

	Blance = 0

	// for i = 0; i < bc.Len; i++ {
	for i = bc.Len - 1; i < bc.Len; i-- {
		cb = (*bc.BodyArr)[i].QueryBlance(&Addr)
		if cb >= 0 {
			Blance = cb
			break
		}
	}

	return Blance
}
