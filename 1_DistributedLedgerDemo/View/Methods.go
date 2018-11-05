package View

import (
	"XGQ/1_DistributedLedgerDemo/Model"
	"fmt"
	"time"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/rtl"
	"github.com/ying32/govcl/vcl/types"

	kis "XGQ/2_KademliaSimu/KademliaInformationSystem"
)

// func (MF *MainForm) Initilization(vci_l *ViewCtrlInterface) {
// vci = vci_l
func (MF *MainForm) Initilization() {
	MF.createCtrl()
	MF.regEvt()
}

func (MF *MainForm) ShowAccountBlance(AccountArr map[string]Model.XgqAccount) {
	var item *vcl.TListItem

	MF.lvAccountBlance.Items().BeginUpdate()
	MF.lvAccountBlance.Items().Clear()

	for key, value := range AccountArr {
		item = MF.lvAccountBlance.Items().Add()

		item.SetCaption(fmt.Sprintf("%v", value.Nickname))
		item.SubItems().Add(key)
		item.SubItems().Add(fmt.Sprintf("%v", value.Balance))
	}

	MF.lvAccountBlance.Items().EndUpdate()
}

func (MF *MainForm) ShowBlockChain(pow Model.ProofOfWork) {
	BC := pow.BC

	var i uint64
	var item *vcl.TListItem
	var Head Model.BlockHead
	var Body *Model.BlockData

	MF.lvBlockChain.Items().BeginUpdate()
	MF.lvBlockChain.Items().Clear()

	item = MF.lvBlockChain.Items().Add()
	item.SetCaption("区块链长度")
	item.SubItems().Add(fmt.Sprintf("%v", BC.Len))

	for i = 0; i < BC.Len; i++ {
		MF.addLvBlockChainGrp(fmt.Sprintf("Block#%v", i))
	}

	for i = 0; i < BC.Len; i++ {
		Head = BC.HeadArr[i]
		Body = &(*BC.BodyArr)[i]

		showBlockHead(Head, MF.lvBlockChain, i, BC.Len)
		showBlockBody(Body, MF.lvBlockChain, i)

	}

	MF.lvBlockChain.Items().EndUpdate()

	MF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":账本刷新完成")
}

func (MF *MainForm) ShowKBuckets(kb *kis.KBuckets) {
	var i, j uint32
	var item *vcl.TListItem
	var bucket kis.Bucket

	MF.lvNode.Items().BeginUpdate()
	MF.lvNode.Items().Clear()

	for i = 0; i < kb.BitSpaceNum; i++ {
		if kb.BucketsCnt[i] < 1 {
			continue
		}

		for j = 0; j < kb.BucketsCnt[i]; j++ {
			bucket = (*kb.Buckets)[i][j]
			item = MF.lvNode.Items().Add()

			item.SetCaption(fmt.Sprintf("%s",
				xu.IP2StrWithoutPort(bucket.ExternalIP)))

			item.SubItems().Add(
				fmt.Sprintf("%d", bucket.ExternalPort))
			item.SubItems().Add(bucket.UUID)
		}
	}

	MF.lvNode.Items().EndUpdate()

	MF.StatusBar.Panels().Items(1).SetText(xu.GetCurTimeStr() + ":列表刷新完成")
}

func (MF *MainForm) addLvBlockChainGrp(caption string) {
	Grp := MF.lvBlockChain.Groups().Add()

	Grp.SetFooterAlign(types.TaCenter)
	Grp.SetHeaderAlign(types.TaCenter)
	Grp.SetHeader(caption)
	//state := Grp.State() //默认为0
	state := types.TListGroupStateSet(rtl.Include(0, types.LgsCollapsible))
	Grp.SetState(state)
	Grp.SetTitleImage(-1)
}

func showBlockHead(
	Head Model.BlockHead, lv *vcl.TListView,
	GrpID uint64, GrpNum uint64) {
	GrpIDInt32 := int32(GrpID)

	item := lv.Items().Add()
	item.SetGroupID(GrpIDInt32)
	item.SetCaption("Index")
	item.SubItems().Add(fmt.Sprintf("%v/%v", Head.Index, GrpNum))

	item = lv.Items().Add()
	item.SetGroupID(GrpIDInt32)

	StrTime := time.Unix(Head.Timestamp, 0).Format("2006-01-02 15:04:05")
	item.SetCaption("时间")
	item.SubItems().Add(StrTime)

	item = lv.Items().Add()
	item.SetGroupID(GrpIDInt32)
	item.SetCaption("PreHash")
	item.SubItems().Add(Head.PreHash)

	item = lv.Items().Add()
	item.SetGroupID(GrpIDInt32)
	item.SetCaption("Hash")
	item.SubItems().Add(Head.Hash)

	item = lv.Items().Add()
	item.SetGroupID(GrpIDInt32)
	item.SetCaption("BodyHash")
	item.SubItems().Add(Head.BodyHash)

	item = lv.Items().Add()
	item.SetGroupID(GrpIDInt32)
	item.SetCaption("Nonce")
	item.SubItems().Add(fmt.Sprintf("%v", Head.Nonce))
}

func showBlockBody(
	Body *Model.BlockData, lv *vcl.TListView, GrpID uint64) {
	GrpIDInt32 := int32(GrpID)

	item := lv.Items().Add()
	item.SetGroupID(GrpIDInt32)
	item.SetCaption("交易数量")
	item.SubItems().Add(fmt.Sprintf("%v", Body.Len))

	var i, j uint64
	var Tx Model.Transaction
	var utxo Model.UTXO

	for i = 0; i < Body.Len; i++ {
		Tx = Body.TxArr[i]

		item = lv.Items().Add()
		item.SetGroupID(GrpIDInt32)
		item.SetCaption(" 交易人数")
		item.SubItems().Add(fmt.Sprintf("%v", Tx.Num))

		item = lv.Items().Add()
		item.SetGroupID(GrpIDInt32)
		item.SetCaption(" Hash")
		item.SubItems().Add(Tx.Hash)

		item = lv.Items().Add()
		item.SetGroupID(GrpIDInt32)
		item.SetCaption(" Signature")
		item.SubItems().Add(Tx.Signature)

		for j = 0; j < Tx.Num; j++ {
			item = lv.Items().Add()
			item.SetGroupID(GrpIDInt32)
			item.SetCaption(fmt.Sprintf("  PreTxHash[%v/%v]", j+1, Tx.Num))
			item.SubItems().Add(Tx.PreTxHash[j])

			utxo = Tx.UTXOArr[j]

			item = lv.Items().Add()
			item.SetGroupID(GrpIDInt32)
			item.SetCaption(fmt.Sprintf("  UTXO Address"))
			item.SubItems().Add(utxo.Address)

			item = lv.Items().Add()
			item.SetGroupID(GrpIDInt32)
			item.SetCaption(fmt.Sprintf("  UTXO Balance"))
			item.SubItems().Add(fmt.Sprintf("%v", utxo.Balance))
		}
	}
}

func getListBoxRowNum(
	ListBox vcl.TListBox, RowNum *int32) {
	var i int32
	*RowNum = 0
	Items := ListBox.Items()
	StrAll := Items.Text()
	if StrAll == "" {
		return
	}
	i = 0

	StrAccu := ""
	for {
		cs := Items.Strings(i)
		StrAccu += cs + "\r\n"

		if StrAccu == StrAll {
			break
		}
		i++
	}

	*RowNum = i + 1
}

func (MF *MainForm) ShowIPAddr(p *kis.Peer) {
	var IPStr string

	MF.EditSeedNodeUUID.SetText(p.SeedNodeUUID)
	MF.LabelUUID.SetCaption("本机UUID:" + p.UUID)

	IPStr = xu.IP2StrWithoutPort(p.InternalIP)
	MF.LabelInternalIP.SetCaption("内网IP:" + IPStr)
	IPStr = xu.IP2StrWithoutPort(p.ExternalIP)
	MF.LabelExternalIP.SetCaption("外网IP:" + IPStr)

	MF.EditInternalPort.SetText(fmt.Sprintf("%d", p.InternalPort))
	MF.EditExternalPort.SetText(fmt.Sprintf("%d", p.ExternalPort))
}
