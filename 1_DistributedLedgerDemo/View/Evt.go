package View

import (
	"XGQ/1_DistributedLedgerDemo/Model"
	"strconv"
	"strings"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

func mainFormExitEvt(Sender vcl.IObject, CanClose *bool) {
	*CanClose =
		vcl.MessageDlg("确定退出？", types.MtConfirmation, types.MbYes, types.MbNo) == types.IdYes
}

func (MF *MainForm) sheetBlockMgrOnShow(sender vcl.IObject) {
	VCI.EvtShowBlockChain()
}

func (MF *MainForm) sheetTxOnShow(sender vcl.IObject) {
	VCI.EvtsheetTxOnShow()
}

func (MF *MainForm) BtnSelectDraweeClick(sender vcl.IObject) {
	item := MF.lvAccountBlance.Selected()

	if item.Instance() == 0 {
		MF.StatusBar.Panels().Items(1).SetText("请选择付款人")
		return
	}

	MF.ListBoxDrawee.Items().Clear()
	SubItems := item.SubItems()
	MF.ListBoxDrawee.Items().Add(
		item.Caption() + "-" + SubItems.Strings(0))
}

func (MF *MainForm) BtnAdd2PayeeClick(sender vcl.IObject) {
	item := MF.lvAccountBlance.Selected()

	if item.Instance() == 0 {
		MF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":请选择收款人")
		return
	}

	if MF.EditPayment.Text() == "" {
		MF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":请填写金额")
		return
	}

	PaymentStr := MF.EditPayment.Text()
	Payment, _ := strconv.ParseFloat(PaymentStr, 64)

	SubItems := item.SubItems()
	MF.ListBoxPayee.Items().Add(
		item.Caption() + "-" + SubItems.Strings(0) +
			"-" + strconv.FormatFloat(Payment, 'f', 3, 64))
}

func (MF *MainForm) BtnSubmitTxClick(sender vcl.IObject) {
	Items := MF.ListBoxDrawee.Items()
	if Items.Text() == "" {
		MF.StatusBar.Panels().Items(1).SetText("请选择付款人")
		return
	}

	DraweeAddr := Items.Strings(0)
	Idx := strings.LastIndex(DraweeAddr, "-")
	if Idx > 0 {
		DraweeAddr = DraweeAddr[Idx+1 : len(DraweeAddr)]
	}

	Items = MF.ListBoxPayee.Items()
	if Items.Text() == "" {
		MF.StatusBar.Panels().Items(1).SetText("请选择收款人")
		return
	}

	var i, PayeeNum int32
	var PayeeAddrArr []string
	var CurPayeeAddr, PaymentStr string
	var PayValueArr []float64

	getListBoxRowNum(*MF.ListBoxPayee, &PayeeNum)
	PayeeAddrArr = make([]string, PayeeNum, PayeeNum)
	PayValueArr = make([]float64, PayeeNum+1, PayeeNum+1)

	for i = 0; i < PayeeNum; i++ {
		CurPayeeAddr = Items.Strings(i)
		Idx = strings.LastIndex(CurPayeeAddr, "-")
		if Idx > 0 {
			PaymentStr = CurPayeeAddr[Idx+1 : len(CurPayeeAddr)]
			CurPayeeAddr = CurPayeeAddr[0:Idx]

			Idx = strings.LastIndex(CurPayeeAddr, "-")
			CurPayeeAddr = CurPayeeAddr[Idx+1 : len(CurPayeeAddr)]
		} else {
			MF.StatusBar.Panels().Items(1).SetText("收款人格式有误")
			return
		}

		PayValueArr[i+1], _ = strconv.ParseFloat(PaymentStr, 64)
		PayeeAddrArr[i] = CurPayeeAddr
	}

	VCI.SubmitTx(DraweeAddr, PayeeAddrArr, PayValueArr)

	MF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":交易提交成功")
}

func (MF *MainForm) BtnDelFromPayeesClick(sender vcl.IObject) {
	MF.ListBoxPayee.DeleteSelected()
}

func (MF *MainForm) BtnCheckBlanceClick(sender vcl.IObject) {
	VCI.EvtsheetTxOnShow()
}

func (MF *MainForm) BtnCreateAccountClick(sender vcl.IObject) {
	XgqA := Model.XgqAccount{}
	XgqA.GenRandAccount()

	VCI.Pksf.XgqA = XgqA
	VCI.Pksf.Form.Show()
}

func (MF *MainForm) BtnMineClick(sender vcl.IObject) {
	MF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":挖矿中...")
	MF.BtnMine.SetEnabled(false)

	VCI.EvtMine()
	VCI.EvtShowBlockChain()

	MF.BtnMine.SetEnabled(true)
	MF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":挖矿完成")
}

func (MF *MainForm) BtnLedgerSyncClick(sender vcl.IObject) {
	MF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":开始同步...")
	MF.BtnLedgerSync.SetEnabled(false)

	VCI.EvtLedgerSync()
	VCI.EvtShowBlockChain()

	MF.BtnLedgerSync.SetEnabled(true)
	MF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":同步完成")
}

func (MF *MainForm) BtnRefreshLVBlockChainClick(sender vcl.IObject) {
	VCI.EvtShowBlockChain()
}

func (f *MainForm) OnFormCreate(sender vcl.IObject) {
}

func (MF *MainForm) BtnConnectSeedNodeClick(sender vcl.IObject) {
	MF.BtnConnectSeedNode.SetEnabled(false)
	VCI.ConnectSeedNode()
	MF.BtnConnectSeedNode.SetEnabled(true)
}

func (MF *MainForm) BtnSaveCfgClick(sender vcl.IObject) {
	VCI.SaveNetWorkCfg()
}

func (MF *MainForm) BtnAddUPnPMappingClick(sender vcl.IObject) {
	VCI.AddUPnPMapping()
}

func (MF *MainForm) BtnDeleteUPnPMappingClick(sender vcl.IObject) {
	VCI.DeleteUPnPMapping()
}

func (MF *MainForm) BtnRefreshLvVNodeClick(sender vcl.IObject) {
	VCI.ShowKBuckets()
}
