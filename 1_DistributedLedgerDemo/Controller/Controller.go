package Controller

import (
	"XGQ/1_DistributedLedgerDemo/Model"
	"XGQ/1_DistributedLedgerDemo/View"
	kis "XGQ/2_KademliaSimu/KademliaInformationSystem"
	"XGQ/XgqUtility"
	"time"

	// "github.com/mattn/go-sqlite3"
	"github.com/ying32/govcl/vcl"
)

var xu XgqUtility.Utility

type Controller struct {
	mF   *View.MainForm
	amgr *AccountMgr

	blockChain *Model.BlockChain
	pow        *Model.ProofOfWork
	peer       *kis.Peer

	vci *View.ViewCtrlInterface
	kci *kis.KisCtrlInterface
}

func (ctrl *Controller) Run() {
	ctrl.Initilization()

	// All form creation go before:vcl.Application.Run()
	ctrl.vci.ChViewInit <- true
	vcl.Application.Run()
}

func (ctrl *Controller) ctrlLop() {
	// time.Sleep(2000 * time.Millisecond)

	// ctrl.amgr.AccountArr[0].Balance = 50

	<-ctrl.vci.ChViewInit
	close(ctrl.vci.ChViewInit)

	for {
		time.Sleep(3000 * time.Millisecond)

		// ctrl.mF.ShowAccountBlance(ctrl.amgr.AccountArr)

	}
}
