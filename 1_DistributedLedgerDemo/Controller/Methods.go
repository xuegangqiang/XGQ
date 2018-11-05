package Controller

import (
	"XGQ/1_DistributedLedgerDemo/Model"
	"fmt"

	kis "XGQ/2_KademliaSimu/KademliaInformationSystem"
)

func (ctrl *Controller) Mine() {
	ctrl.pow.Mine()
}

func (ctrl *Controller) LedgerSync() {
	ctrl.peer.CallRPCBlockChainSync(
		ctrl.mF.EditSeedNodeAddr.Text())

	ctrl.ShowBlockChain()
}

func (ctrl *Controller) SubmitTx(
	DraweeAddr string,
	PayeeAddrArr []string,
	PayValueArr []float64) {

	amgr := ctrl.amgr
	pow := ctrl.pow

	DraweeAccount := amgr.AccountArr[DraweeAddr]
	PayeeNum := len(PayeeAddrArr)
	TxNum := PayeeNum + 1

	AddrArr := make([]string, TxNum, TxNum)
	LastValueArr := make([]float64, TxNum, TxNum)
	UtxoHashArr := make([]string, TxNum, TxNum)

	AddrArr[0] = DraweeAccount.PublicKeyHash
	LastValueArr[0] = DraweeAccount.Balance
	PayValueArr[0] = 0
	UtxoHashArr[0] = DraweeAccount.UtxoHash

	var PayeeAccount Model.XgqAccount
	for i := 1; i < TxNum; i++ {
		PayeeAccount = amgr.AccountArr[PayeeAddrArr[i-1]]
		AddrArr[i] = PayeeAccount.PublicKeyHash
		LastValueArr[i] = PayeeAccount.Balance
		UtxoHashArr[i] = PayeeAccount.UtxoHash
	}

	Tx, err := DraweeAccount.Pay(
		AddrArr,
		LastValueArr,
		PayValueArr,
		UtxoHashArr)
	if err != nil {
		ctrl.mF.StatusBar.Panels().Items(1).SetText("交易提交失败:付款人余额不足")
	}

	pow.Add2TrasactionPool(Tx, &amgr.AccountArr)
}

func (ctrl *Controller) QueryBlance() {
	amgr := ctrl.amgr
	BC := ctrl.pow.BC
	XgqA := Model.XgqAccount{}

	for key, _ := range amgr.AccountArr {
		XgqA = amgr.AccountArr[key]
		XgqA.Balance = BC.QueryBlance(key)
		amgr.AccountArr[key] = XgqA
	}

	ctrl.mF.ShowAccountBlance(ctrl.amgr.AccountArr)

	ctrl.mF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":余额查询成功")
}

func (ctrl *Controller) AddAccount(XgqA *Model.XgqAccount) {
	ctrl.amgr.AccountArr[XgqA.PublicKeyHash] = *XgqA
	ctrl.QueryBlance()
}

func (ctrl *Controller) ShowBlockChain() {
	ctrl.mF.ShowBlockChain(*ctrl.pow)
}

func (ctrl *Controller) ConnectSeedNode() {
	var RefNodeArr []kis.Peer

	ctrl.mF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":尝试连接种子节点...")
	RefNodeArr = make([]kis.Peer, 1, 1)
	RefNodeArr[0] = *(kis.NewPeer())

	RefNodeArr[0].UUID = ctrl.mF.EditSeedNodeUUID.Text()
	RefNodeArr[0].ExternalIP, RefNodeArr[0].ExternalPort =
		xu.Str2IP(ctrl.mF.EditSeedNodeAddr.Text())

	ctrl.peer.JoinNetWork(&RefNodeArr)

	// Test RPCFindNode method
	// ctrl.peer.CallRPCFindNode(
	// 	ctrl.mF.EditSeedNodeAddr.Text(),
	// 	ctrl.mF.EditSeedNodeUUID.Text())

	// Test RPCPing method
	// DstAddr := ctrl.mF.EditSeedNodeAddr.Text()
	// ctrl.peer.CallRPCPing(DstAddr)

	kb := ctrl.peer.GetKBuckets()
	// kb.AddBucket(*RefNodeArr[0].FormABucket())

	ctrl.mF.ShowKBuckets(kb)
	ctrl.mF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":种子节点连接完成")
}

func (ctrl *Controller) SaveNetWorkCfg() {
	var Str, info string
	var err error

	ctrl.mF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":尝试保存网络参数...")

	Str = ctrl.mF.EditSeedNodeAddr.Text()
	ctrl.peer.SeedNodeIP, ctrl.peer.SeedNodePort = xu.Str2IP(Str)
	ctrl.peer.SeedNodeUUID = ctrl.mF.EditSeedNodeUUID.Text()

	fmt.Sscanf(ctrl.mF.EditInternalPort.Text(), "%d", &ctrl.peer.InternalPort)
	fmt.Sscanf(ctrl.mF.EditExternalPort.Text(), "%d", &ctrl.peer.ExternalPort)

	Str = ctrl.mF.LabelInternalIP.Caption()
	ctrl.peer.InternalIP, _ = xu.Str2IP(Str[9:])
	Str = ctrl.mF.LabelExternalIP.Caption()
	ctrl.peer.ExternalIP, _ = xu.Str2IP(Str[9:])

	err = ctrl.peer.SaveCfgFile("")
	if err != nil {
		info = fmt.Sprintf("保存网络参数失败:%s", err)
		ctrl.mF.StatusBar.Panels().Items(1).SetText(info)
		return
	}

	ctrl.mF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":保存网络参数成功")
}

func (ctrl *Controller) AddUPnPMapping() {
	var nwp kis.NetWorkProtocol
	var InternalPort, ExternalPort int
	nwp = kis.NetWorkProtocol{}

	ctrl.mF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":尝试增加映射...")
	fmt.Sscanf(ctrl.mF.EditInternalPort.Text(), "%d", &InternalPort)
	fmt.Sscanf(ctrl.mF.EditExternalPort.Text(), "%d", &ExternalPort)
	if !nwp.AddUPnPMapping(InternalPort, ExternalPort) {
		ctrl.mF.StatusBar.Panels().Items(1).SetText("增加映射失败")
		return
	}
	ctrl.mF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":增加映射成功")
}

func (ctrl *Controller) DeleteUPnPMapping() {
	var nwp kis.NetWorkProtocol
	var InternalPort, ExternalPort int
	nwp = kis.NetWorkProtocol{}

	ctrl.mF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":尝试删除映射...")
	fmt.Sscanf(ctrl.mF.EditInternalPort.Text(), "%d", &InternalPort)
	fmt.Sscanf(ctrl.mF.EditExternalPort.Text(), "%d", &ExternalPort)
	if !nwp.DeleteUPnPMapping(InternalPort, ExternalPort) {
		ctrl.mF.StatusBar.Panels().Items(1).SetText("删除映射失败")
		return
	}
	ctrl.mF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":删除映射成功")
}

func (ctrl *Controller) ShowKBuckets() {
	ctrl.mF.ShowKBuckets(ctrl.peer.GetKBuckets())
}

// func (ctrl *Controller) NewKisCtrlInterface() *kis.KisCtrlInterface {
// 	var kci kis.KisCtrlInterface
// 	kci = kis.KisCtrlInterface{}

// 	return &kci
// }

func (ctrl *Controller) Report2Ctrl(msg string) {
	ctrl.mF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":" + msg)
}

func (ctrl *Controller) GetPOW() *Model.ProofOfWork {
	return ctrl.pow
}
