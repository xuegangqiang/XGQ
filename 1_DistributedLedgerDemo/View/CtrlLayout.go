package View

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/rtl"
	"github.com/ying32/govcl/vcl/types"
)

func (MF *MainForm) createCtrl() {
	vcl.Application.SetIconResId(3)
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)

	MF.layoutCtrl()
}

func (MF *MainForm) layoutCtrl() {
	layoutMainForm(MF)

	layoutSheetBlockMgr(MF)
	layoutSheetTx(MF)
	layoutSheetNet(MF)
}

func layoutMainForm(MF *MainForm) {
	// mainForm
	MF.mainForm = vcl.Application.CreateForm()
	MF.mainForm.SetCaption("分布式账本示例")
	MF.mainForm.SetPosition(types.PoScreenCenter)
	// MF.mainForm.EnabledMaximize(false)
	MF.mainForm.SetWidth(900)
	MF.mainForm.SetHeight(600)
	MF.mainForm.SetDoubleBuffered(true)

	// pageCtrl
	MF.pageCtrl = vcl.NewPageControl(MF.mainForm)
	MF.pageCtrl.SetParent(MF.mainForm)
	MF.pageCtrl.SetAlign(types.AlClient)

	// sheetBlockMgr
	MF.sheetBlockMgr = vcl.NewTabSheet(MF.mainForm)
	MF.sheetBlockMgr.SetPageControl(MF.pageCtrl)
	MF.sheetBlockMgr.SetCaption("挖矿")

	// sheetTx
	MF.sheetTx = vcl.NewTabSheet(MF.mainForm)
	MF.sheetTx.SetPageControl(MF.pageCtrl)
	MF.sheetTx.SetCaption("交易")

	// sheetNet
	MF.sheetNet = vcl.NewTabSheet(MF.mainForm)
	MF.sheetNet.SetPageControl(MF.pageCtrl)
	MF.sheetNet.SetCaption("网络")

	// StatusBar
	MF.StatusBar = vcl.NewStatusBar(MF.mainForm)
	MF.StatusBar.SetParent(MF.mainForm)

	spnl := MF.StatusBar.Panels().Add()
	spnl.SetText("cryptoxgq@163.com")
	spnl.SetWidth(160)

	spnl = MF.StatusBar.Panels().Add()
	spnl.SetText("")
	spnl.SetWidth(400)
}

func layoutSheetBlockMgr(MF *MainForm) {
	// lvBlockChain
	MF.lvBlockChain = vcl.NewListView(MF.sheetBlockMgr)
	MF.lvBlockChain.SetParent(MF.sheetBlockMgr)
	MF.lvBlockChain.SetLeft(0)
	MF.lvBlockChain.SetTop(0)
	MF.lvBlockChain.SetWidth(820)
	MF.lvBlockChain.SetHeight(520)
	AkBlockChain := types.TAnchors(
		rtl.Include(0, types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))
	MF.lvBlockChain.SetAnchors(AkBlockChain)
	MF.lvBlockChain.SetRowSelect(true)
	MF.lvBlockChain.SetReadOnly(true)
	MF.lvBlockChain.SetViewStyle(types.VsReport)
	MF.lvBlockChain.SetGridLines(true)
	MF.lvBlockChain.SetColumnClick(false)
	MF.lvBlockChain.SetGroupView(true)
	col := MF.lvBlockChain.Columns().Add()
	col.SetCaption("项目")
	col.SetWidth(200)
	col = MF.lvBlockChain.Columns().Add()
	col.SetCaption("值")
	col.SetWidth(600)

	// BtnMine
	MF.BtnMine = vcl.NewButton(MF.sheetBlockMgr)
	MF.BtnMine.SetParent(MF.sheetBlockMgr)
	MF.BtnMine.SetCaption("挖矿")
	Ak := types.TAnchors(rtl.Include(0, types.AkRight, types.AkTop))
	MF.BtnMine.SetAnchors(Ak)
	MF.BtnMine.SetBounds(822, 5, 50, 30)
	MF.BtnMine.SetOnClick(MF.BtnMineClick)

	// BtnLedgerSync
	MF.BtnLedgerSync = vcl.NewButton(MF.sheetBlockMgr)
	MF.BtnLedgerSync.SetParent(MF.sheetBlockMgr)
	MF.BtnLedgerSync.SetCaption("同步")
	Ak = types.TAnchors(rtl.Include(0, types.AkRight, types.AkTop))
	MF.BtnLedgerSync.SetAnchors(Ak)
	MF.BtnLedgerSync.SetBounds(822, 60, 50, 30)
	MF.BtnLedgerSync.SetOnClick(MF.BtnLedgerSyncClick)

	// BtnRefreshLVBlockChain
	MF.BtnRefreshLVBlockChain = vcl.NewButton(MF.sheetBlockMgr)
	MF.BtnRefreshLVBlockChain.SetParent(MF.sheetBlockMgr)
	MF.BtnRefreshLVBlockChain.SetCaption("刷新")
	Ak = types.TAnchors(rtl.Include(0, types.AkRight, types.AkTop))
	MF.BtnRefreshLVBlockChain.SetAnchors(Ak)
	MF.BtnRefreshLVBlockChain.SetBounds(822, 115, 50, 30)
	MF.BtnRefreshLVBlockChain.SetOnClick(MF.BtnRefreshLVBlockChainClick)
}

func layoutSheetTx(MF *MainForm) {
	// lvAccountBlance
	MF.lvAccountBlance = vcl.NewListView(MF.sheetTx)
	MF.lvAccountBlance.SetParent(MF.sheetTx)
	MF.lvAccountBlance.SetLeft(0)
	MF.lvAccountBlance.SetTop(0)
	MF.lvAccountBlance.SetWidth(600)
	MF.lvAccountBlance.SetHeight(100)
	MF.lvAccountBlance.SetRowSelect(true)
	MF.lvAccountBlance.SetReadOnly(true)
	MF.lvAccountBlance.SetViewStyle(types.VsReport)
	MF.lvAccountBlance.SetGridLines(true)
	// MF.lvAccountBlance.SetColumnClick(false)
	col := MF.lvAccountBlance.Columns().Add()
	col.SetCaption("昵称")
	col.SetWidth(50)
	col = MF.lvAccountBlance.Columns().Add()
	col.SetCaption("地址")
	col.SetWidth(400)
	col = MF.lvAccountBlance.Columns().Add()
	col.SetCaption("余额")
	col.SetWidth(50)

	// BtnSelectDrawee
	MF.BtnSelectDrawee = vcl.NewButton(MF.sheetTx)
	MF.BtnSelectDrawee.SetParent(MF.sheetTx)
	MF.BtnSelectDrawee.SetCaption("选择付款人")
	Ak := types.TAnchors(rtl.Include(0, types.AkLeft, types.AkTop))
	MF.BtnSelectDrawee.SetAnchors(Ak)
	MF.BtnSelectDrawee.SetBounds(650, 5, 90, 30)
	MF.BtnSelectDrawee.SetOnClick(MF.BtnSelectDraweeClick)

	// BtnCheckBlance
	MF.BtnCheckBlance = vcl.NewButton(MF.sheetTx)
	MF.BtnCheckBlance.SetParent(MF.sheetTx)
	MF.BtnCheckBlance.SetCaption("余额查询")
	MF.BtnCheckBlance.SetAnchors(Ak)
	MF.BtnCheckBlance.SetBounds(750, 5, 90, 30)
	MF.BtnCheckBlance.SetOnClick(MF.BtnCheckBlanceClick)

	// BtnAdd2Payee
	MF.BtnAdd2Payee = vcl.NewButton(MF.sheetTx)
	MF.BtnAdd2Payee.SetParent(MF.sheetTx)
	MF.BtnAdd2Payee.SetCaption("添加收款人")
	MF.BtnAdd2Payee.SetAnchors(Ak)
	MF.BtnAdd2Payee.SetBounds(650, 40, 90, 30)
	MF.BtnAdd2Payee.SetOnClick(MF.BtnAdd2PayeeClick)

	// BtnCreateAccount
	MF.BtnCreateAccount = vcl.NewButton(MF.sheetTx)
	MF.BtnCreateAccount.SetParent(MF.sheetTx)
	MF.BtnCreateAccount.SetCaption("新建账户")
	MF.BtnCreateAccount.SetAnchors(Ak)
	MF.BtnCreateAccount.SetBounds(750, 40, 90, 30)
	MF.BtnCreateAccount.SetOnClick(MF.BtnCreateAccountClick)

	// LabelDrawee
	MF.LabelDrawee = vcl.NewLabel(MF.sheetTx)
	MF.LabelDrawee.SetParent(MF.sheetTx)
	Ak = types.TAnchors(rtl.Include(0, types.AkLeft, types.AkTop))
	MF.LabelDrawee.SetAnchors(Ak)
	MF.LabelDrawee.SetBounds(0, 130, 60, 30)
	MF.LabelDrawee.SetCaption("付款人")

	// LabelPayee
	MF.LabelPayee = vcl.NewLabel(MF.sheetTx)
	MF.LabelPayee.SetParent(MF.sheetTx)
	MF.LabelPayee.SetAnchors(Ak)
	MF.LabelPayee.SetBounds(0, 230, 60, 30)
	MF.LabelPayee.SetCaption("收款人")

	// LabelPayment
	MF.LabelPayment = vcl.NewLabel(MF.sheetTx)
	MF.LabelPayment.SetParent(MF.sheetTx)
	MF.LabelPayment.SetAnchors(Ak)
	MF.LabelPayment.SetBounds(650, 80, 40, 30)
	MF.LabelPayment.SetCaption("金额：")

	// ListBoxDrawee
	MF.ListBoxDrawee = vcl.NewListBox(MF.sheetTx)
	MF.ListBoxDrawee.SetParent(MF.sheetTx)
	MF.ListBoxDrawee.SetAnchors(Ak)
	MF.ListBoxDrawee.SetBounds(0, 150, 500, 30)

	// ListBoxPayee
	MF.ListBoxPayee = vcl.NewListBox(MF.sheetTx)
	MF.ListBoxPayee.SetParent(MF.sheetTx)
	MF.ListBoxPayee.SetAnchors(Ak)
	MF.ListBoxPayee.SetBounds(0, 250, 500, 80)

	// BtnSubmitTx
	MF.BtnSubmitTx = vcl.NewButton(MF.sheetTx)
	MF.BtnSubmitTx.SetParent(MF.sheetTx)
	MF.BtnSubmitTx.SetCaption("提交交易")
	MF.BtnSubmitTx.SetAnchors(Ak)
	MF.BtnSubmitTx.SetBounds(530, 150, 90, 30)
	MF.BtnSubmitTx.SetOnClick(MF.BtnSubmitTxClick)

	// BtnDelFromPayees
	MF.BtnDelFromPayees = vcl.NewButton(MF.sheetTx)
	MF.BtnDelFromPayees.SetParent(MF.sheetTx)
	MF.BtnDelFromPayees.SetCaption("删除收款人")
	MF.BtnDelFromPayees.SetAnchors(Ak)
	MF.BtnDelFromPayees.SetBounds(530, 250, 90, 30)
	MF.BtnDelFromPayees.SetOnClick(MF.BtnDelFromPayeesClick)

	// EditPayment
	MF.EditPayment = vcl.NewEdit(MF.sheetTx)
	MF.EditPayment.SetParent(MF.sheetTx)
	MF.EditPayment.SetAnchors(Ak)
	MF.EditPayment.SetBounds(690, 75, 45, 30)
}

func layoutSheetNet(MF *MainForm) {
	Ak := types.TAnchors(rtl.Include(0, types.AkLeft, types.AkTop))

	// EditSeedNodeAddr
	MF.EditSeedNodeAddr = vcl.NewEdit(MF.sheetNet)
	MF.EditSeedNodeAddr.SetParent(MF.sheetNet)
	MF.EditSeedNodeAddr.SetAnchors(Ak)
	MF.EditSeedNodeAddr.SetBounds(85, 15-2, 130, 30)
	MF.EditSeedNodeAddr.SetText("192.168.1.4:56068")

	// EditSeedNodeUUID
	MF.EditSeedNodeUUID = vcl.NewEdit(MF.sheetNet)
	MF.EditSeedNodeUUID.SetParent(MF.sheetNet)
	MF.EditSeedNodeUUID.SetAnchors(Ak)
	MF.EditSeedNodeUUID.SetBounds(20, 15-2+25, 220, 30)
	MF.EditSeedNodeUUID.SetText("SeedNodeUUID")

	// LabelSeedNode
	MF.LabelSeedNode = vcl.NewLabel(MF.sheetNet)
	MF.LabelSeedNode.SetParent(MF.sheetNet)
	MF.LabelSeedNode.SetAnchors(Ak)
	MF.LabelSeedNode.SetBounds(20, 15, 60, 30)
	MF.LabelSeedNode.SetCaption("种子节点:")

	// LabelUUID
	MF.LabelUUID = vcl.NewLabel(MF.sheetNet)
	MF.LabelUUID.SetParent(MF.sheetNet)
	MF.LabelUUID.SetAnchors(Ak)
	MF.LabelUUID.SetBounds(20, 50+5+15, 200, 30)
	MF.LabelUUID.SetCaption("本机UUID")

	// LabelInternalIP
	MF.LabelInternalIP = vcl.NewLabel(MF.sheetNet)
	MF.LabelInternalIP.SetParent(MF.sheetNet)
	MF.LabelInternalIP.SetAnchors(Ak)
	MF.LabelInternalIP.SetBounds(20, 50+2+50, 120, 30)
	MF.LabelInternalIP.SetCaption("内网IP:192.168.0.xx")

	// LabelExternalIP
	MF.LabelExternalIP = vcl.NewLabel(MF.sheetNet)
	MF.LabelExternalIP.SetParent(MF.sheetNet)
	MF.LabelExternalIP.SetAnchors(Ak)
	MF.LabelExternalIP.SetBounds(20, 90+50, 120, 30)
	MF.LabelExternalIP.SetCaption("公网IP:192.168.1.xx")

	// EditInternalPort
	MF.EditInternalPort = vcl.NewEdit(MF.sheetNet)
	MF.EditInternalPort.SetParent(MF.sheetNet)
	MF.EditInternalPort.SetAnchors(Ak)
	MF.EditInternalPort.SetBounds(145, 50+50, 45, 30)
	MF.EditInternalPort.SetText("11200")

	// EditIExternalPort
	MF.EditExternalPort = vcl.NewEdit(MF.sheetNet)
	MF.EditExternalPort.SetParent(MF.sheetNet)
	MF.EditExternalPort.SetAnchors(Ak)
	MF.EditExternalPort.SetBounds(145, 90-2+50, 45, 30)
	MF.EditExternalPort.SetText("56067")

	// BtnConnectSeedNode
	MF.BtnConnectSeedNode = vcl.NewButton(MF.sheetNet)
	MF.BtnConnectSeedNode.SetParent(MF.sheetNet)
	MF.BtnConnectSeedNode.SetCaption("连接")
	MF.BtnConnectSeedNode.SetAnchors(Ak)
	MF.BtnConnectSeedNode.SetBounds(250, 10+10, 60, 30)
	MF.BtnConnectSeedNode.SetOnClick(MF.BtnConnectSeedNodeClick)

	// BtnSaveCfg
	MF.BtnSaveCfg = vcl.NewButton(MF.sheetNet)
	MF.BtnSaveCfg.SetParent(MF.sheetNet)
	MF.BtnSaveCfg.SetCaption("保存")
	MF.BtnSaveCfg.SetAnchors(Ak)
	MF.BtnSaveCfg.SetBounds(320, 10+10, 40, 30)
	MF.BtnSaveCfg.SetOnClick(MF.BtnSaveCfgClick)

	// BtnAddUPnPMapping
	MF.BtnAddUPnPMapping = vcl.NewButton(MF.sheetNet)
	MF.BtnAddUPnPMapping.SetParent(MF.sheetNet)
	MF.BtnAddUPnPMapping.SetCaption("增加映射")
	MF.BtnAddUPnPMapping.SetAnchors(Ak)
	MF.BtnAddUPnPMapping.SetBounds(220, 45+1+50, 60, 30)
	MF.BtnAddUPnPMapping.SetOnClick(MF.BtnAddUPnPMappingClick)

	// BtnDeleteUPnPMapping
	MF.BtnDeleteUPnPMapping = vcl.NewButton(MF.sheetNet)
	MF.BtnDeleteUPnPMapping.SetParent(MF.sheetNet)
	MF.BtnDeleteUPnPMapping.SetCaption("删除映射")
	MF.BtnDeleteUPnPMapping.SetAnchors(Ak)
	MF.BtnDeleteUPnPMapping.SetBounds(220, 85+50-1, 60, 30)
	MF.BtnDeleteUPnPMapping.SetOnClick(MF.BtnDeleteUPnPMappingClick)

	// LabellvNode
	MF.LabellvNode = vcl.NewLabel(MF.sheetNet)
	MF.LabellvNode.SetParent(MF.sheetNet)
	MF.LabellvNode.SetAnchors(Ak)
	MF.LabellvNode.SetBounds(5, 130+50, 90, 30)
	MF.LabellvNode.SetCaption("已发现的节点:")

	// lvNode
	MF.lvNode = vcl.NewListView(MF.sheetNet)
	MF.lvNode.SetParent(MF.sheetNet)
	MF.lvNode.SetLeft(0)
	MF.lvNode.SetTop(150 + 50)
	MF.lvNode.SetWidth(400)
	MF.lvNode.SetHeight(300)
	MF.lvNode.SetRowSelect(true)
	MF.lvNode.SetReadOnly(true)
	MF.lvNode.SetViewStyle(types.VsReport)
	MF.lvNode.SetGridLines(true)
	col := MF.lvNode.Columns().Add()
	col.SetAlignment(types.TaCenter)
	col.SetCaption("地址")
	col.SetWidth(80)
	col = MF.lvNode.Columns().Add()
	col.SetAlignment(types.TaCenter)
	col.SetCaption("端口")
	col.SetWidth(50)
	col = MF.lvNode.Columns().Add()
	col.SetAlignment(types.TaCenter)
	col.SetCaption("UUID")
	col.SetWidth(220)

	// BtnRefreshLvVNode
	MF.BtnRefreshLvVNode = vcl.NewButton(MF.sheetNet)
	MF.BtnRefreshLvVNode.SetParent(MF.sheetNet)
	MF.BtnRefreshLvVNode.SetCaption("刷新列表")
	MF.BtnRefreshLvVNode.SetAnchors(Ak)
	MF.BtnRefreshLvVNode.SetBounds(410, 200, 60, 30)
	MF.BtnRefreshLvVNode.SetOnClick(MF.BtnRefreshLvVNodeClick)
}
