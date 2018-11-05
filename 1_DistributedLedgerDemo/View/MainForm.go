package View

import (
	"XGQ/1_DistributedLedgerDemo/Model"
	"XGQ/XgqUtility"

	"github.com/ying32/govcl/vcl"
)

var xu XgqUtility.Utility

var VCI *ViewCtrlInterface

type ViewCtrlInterface struct {
	ChView     chan bool
	ChViewInit chan bool

	EvtMine           func()
	EvtLedgerSync     func()
	EvtsheetTxOnShow  func()
	EvtShowBlockChain func()
	SubmitTx          func(string, []string, []float64)
	AddAccount        func(*Model.XgqAccount)

	ConnectSeedNode   func()
	SaveNetWorkCfg    func()
	AddUPnPMapping    func()
	DeleteUPnPMapping func()
	ShowKBuckets      func()

	Pksf *PrivateKeySaveForm
}

// Test result shows that vcl do not allow any concurency touch it
// Otherwise, the main thread will be blocked
// Thus, ALL UI CHANGE SHOULD BE DONE IN THE MAIN THREAD!

type MainForm struct {
	MFPieceForm
	MFPieceSheetBlockMgr
	MFPieceSheetTx
	MFPieceSheetNet
}

type MFPieceForm struct {
	mainForm      *vcl.TForm
	pageCtrl      *vcl.TPageControl
	sheetBlockMgr *vcl.TTabSheet
	sheetTx       *vcl.TTabSheet
	sheetNet      *vcl.TTabSheet
	StatusBar     *vcl.TStatusBar
}

type MFPieceSheetBlockMgr struct {
	BtnMine                *vcl.TButton
	BtnLedgerSync          *vcl.TButton
	BtnRefreshLVBlockChain *vcl.TButton
	lvBlockChain           *vcl.TListView
}

type MFPieceSheetTx struct {
	lvAccountBlance  *vcl.TListView
	BtnSelectDrawee  *vcl.TButton
	BtnAdd2Payee     *vcl.TButton
	BtnSubmitTx      *vcl.TButton
	BtnDelFromPayees *vcl.TButton
	BtnCheckBlance   *vcl.TButton
	BtnCreateAccount *vcl.TButton
	LabelDrawee      *vcl.TLabel
	LabelPayee       *vcl.TLabel
	LabelPayment     *vcl.TLabel
	ListBoxDrawee    *vcl.TListBox
	ListBoxPayee     *vcl.TListBox
	EditPayment      *vcl.TEdit
}

type MFPieceSheetNet struct {
	EditSeedNodeAddr     *vcl.TEdit
	EditSeedNodeUUID     *vcl.TEdit
	EditInternalPort     *vcl.TEdit
	EditExternalPort     *vcl.TEdit
	LabelSeedNode        *vcl.TLabel
	LabelUUID            *vcl.TLabel
	LabelInternalIP      *vcl.TLabel
	LabelExternalIP      *vcl.TLabel
	lvNode               *vcl.TListView
	LabellvNode          *vcl.TLabel
	BtnConnectSeedNode   *vcl.TButton
	BtnAddUPnPMapping    *vcl.TButton
	BtnDeleteUPnPMapping *vcl.TButton
	BtnSaveCfg           *vcl.TButton
	BtnRefreshLvVNode    *vcl.TButton
}
