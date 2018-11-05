package Controller

import (
	"XGQ/1_DistributedLedgerDemo/Model"
	"XGQ/1_DistributedLedgerDemo/View"
	"math/big"
)

type AccountMgr struct {
	AccountArr map[string]Model.XgqAccount
	pksf       View.PrivateKeySaveForm
}

func (amgr *AccountMgr) Initilization() {
	amgr.AccountArr = make(map[string]Model.XgqAccount)
	amgr.pksf = View.PrivateKeySaveForm{}
	amgr.pksf.Initilization()

	fixAccount(amgr)
	setCustomBalance(amgr)
}

func (amgr *AccountMgr) AddAccount(Nickname string, pksf *View.PrivateKeySaveForm) {
	xgqA := Model.XgqAccount{}
	xgqA.GenRandAccount()
	xgqA.Nickname = Nickname

	amgr.AccountArr[xgqA.PublicKeyHash] = xgqA
}

func str2BigInt(str string) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(str, 10)
	return n
}
