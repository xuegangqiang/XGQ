package Controller

import (
	"XGQ/1_DistributedLedgerDemo/Model"
	"XGQ/1_DistributedLedgerDemo/View"
	kis "XGQ/2_KademliaSimu/KademliaInformationSystem"
)

func (ctrl *Controller) Initilization() {
	ctrl.mF = &View.MainForm{}

	// ctrl.InitializeKisCtrlInterface()

	ctrl.pow = &Model.ProofOfWork{}
	ctrl.pow.Initilization()
	ctrl.blockChain = ctrl.pow.BC

	ctrl.InitializeViewCtrlInterface()
	ctrl.mF.Initilization()

	ctrl.amgr = &AccountMgr{}
	ctrl.amgr.Initilization()
	ctrl.vci.Pksf = &ctrl.amgr.pksf

	ctrl.InitializeNetWork()
	addOneBlock(ctrl)

	go ctrl.ctrlLop()
}

func (ctrl *Controller) InitializeViewCtrlInterface() {
	ctrl.vci = &View.ViewCtrlInterface{}

	ctrl.vci.ChViewInit = make(chan bool)
	ctrl.vci.EvtShowBlockChain = ctrl.ShowBlockChain
	ctrl.vci.EvtsheetTxOnShow = ctrl.QueryBlance
	ctrl.vci.EvtMine = ctrl.Mine
	ctrl.vci.EvtLedgerSync = ctrl.LedgerSync
	ctrl.vci.SubmitTx = ctrl.SubmitTx
	ctrl.vci.AddAccount = ctrl.AddAccount

	ctrl.vci.ConnectSeedNode = ctrl.ConnectSeedNode
	ctrl.vci.SaveNetWorkCfg = ctrl.SaveNetWorkCfg
	ctrl.vci.AddUPnPMapping = ctrl.AddUPnPMapping
	ctrl.vci.DeleteUPnPMapping = ctrl.DeleteUPnPMapping
	ctrl.vci.ShowKBuckets = ctrl.ShowKBuckets

	View.VCI = ctrl.vci
	// ctrl.vci.Pksf = &ctrl.amgr.pksf
}

func (ctrl *Controller) InitializeNetWork() {
	ctrl.peer = kis.NewPeer()
	ctrl.peer.Initialization()

	ctrl.mF.ShowIPAddr(ctrl.peer)

	// ctrl.kci = ctrl.NewKisCtrlInterface()
	kis.Kad = kis.NewKademlia()
	kis.Kad.Scenario = kis.RealScenario
	ctrl.InitializeKisCtrlInterface()

	go ctrl.peer.RunRPCServer()
}

func (ctrl *Controller) InitializeKisCtrlInterface() {
	ctrl.kci = &kis.KisCtrlInterface{}
	ctrl.kci.Report2Ctrl = ctrl.Report2Ctrl
	ctrl.kci.GetPOW = ctrl.GetPOW
	kis.KCI = ctrl.kci
}
