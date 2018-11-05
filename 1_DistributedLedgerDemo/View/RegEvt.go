package View

func (MF *MainForm) regEvt() {
	MF.sheetBlockMgr.SetOnShow(MF.sheetBlockMgrOnShow)
	MF.sheetTx.SetOnShow(MF.sheetTxOnShow)
	MF.mainForm.SetOnCloseQuery(mainFormExitEvt)
}
