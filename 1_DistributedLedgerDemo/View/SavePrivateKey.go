package View

import (
	"XGQ/1_DistributedLedgerDemo/Model"
	"crypto/sha512"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/rtl"
	"github.com/ying32/govcl/vcl/types"
)

type PrivateKeySaveForm struct {
	Form          *vcl.TForm
	LabelPassword *vcl.TLabel
	LabelNickname *vcl.TLabel
	EditPassword  *vcl.TEdit
	EditNickname  *vcl.TEdit
	BtnSave       *vcl.TButton
	SaveFileDl    *vcl.TSaveDialog

	addr      string
	x, y, d   big.Int
	curveType int

	sc Model.SealedCrypto

	XgqA Model.XgqAccount
}

func (pksf *PrivateKeySaveForm) SetPrivateKeyData(
	Addr string, X, Y, D big.Int, CurveType int) {
	pksf.addr = Addr
	pksf.x = X
	pksf.y = Y
	pksf.d = D
	pksf.curveType = CurveType
}

func (pksf *PrivateKeySaveForm) Initilization() {
	pksf.sc = Model.SealedCrypto{}

	// Form
	pksf.Form = vcl.Application.CreateForm()
	pksf.Form.SetCaption("新建账户")
	pksf.Form.SetPosition(types.PoScreenCenter)
	pksf.Form.SetWidth(400)
	pksf.Form.SetHeight(200)
	pksf.Form.SetDoubleBuffered(true)

	pksf.layoutCtrl()
}

func (pksf *PrivateKeySaveForm) layoutCtrl() {
	// LabelPassword
	pksf.LabelPassword = vcl.NewLabel(pksf.Form)
	pksf.LabelPassword.SetParent(pksf.Form)
	Ak := types.TAnchors(rtl.Include(0, types.AkLeft, types.AkTop))
	pksf.LabelPassword.SetAnchors(Ak)
	pksf.LabelPassword.SetBounds(5, 10, 90, 30)
	pksf.LabelPassword.SetCaption("私钥加密口令")

	// LabelNickname
	pksf.LabelNickname = vcl.NewLabel(pksf.Form)
	pksf.LabelNickname.SetParent(pksf.Form)
	pksf.LabelNickname.SetAnchors(Ak)
	pksf.LabelNickname.SetBounds(5, 45, 90, 30)
	pksf.LabelNickname.SetCaption("昵称")

	// EditPassword
	pksf.EditPassword = vcl.NewEdit(pksf.Form)
	pksf.EditPassword.SetParent(pksf.Form)
	pksf.EditPassword.SetAnchors(Ak)
	pksf.EditPassword.SetBounds(100, 5, 200, 30)

	// EditNickname
	pksf.EditNickname = vcl.NewEdit(pksf.Form)
	pksf.EditNickname.SetParent(pksf.Form)
	pksf.EditNickname.SetAnchors(Ak)
	pksf.EditNickname.SetBounds(100, 40, 200, 30)

	// BtnSave
	pksf.BtnSave = vcl.NewButton(pksf.Form)
	pksf.BtnSave.SetParent(pksf.Form)
	pksf.BtnSave.SetCaption("保存")
	Ak = types.TAnchors(rtl.Include(0, types.AkRight, types.AkBottom))
	pksf.BtnSave.SetAnchors(Ak)
	pksf.BtnSave.SetBounds(300, 100, 50, 30)
	pksf.BtnSave.SetOnClick(pksf.BtnSaveClick)

	// SaveFileDl
	pksf.SaveFileDl = vcl.NewSaveDialog(pksf.Form)
	pksf.SaveFileDl.SetFilter("DAT文件(*.dat)|*.dat|所有文件(*.*)|*.*")
	pksf.SaveFileDl.SetOptions(
		types.TOpenOptions(rtl.Include(
			uint32(pksf.SaveFileDl.Options()), types.OfShowHelp)))
	pksf.SaveFileDl.SetTitle("保存")
}

func (pksf *PrivateKeySaveForm) BtnSaveClick(sender vcl.IObject) {
	Form := pksf.Form
	SaveFileDl := pksf.SaveFileDl

	SaveFileDl.SetFileName(
		pksf.EditNickname.Text() + "-口令" +
			pksf.EditPassword.Text() + ".dat")
	DlRes := SaveFileDl.Execute()
	if DlRes {
		PrivateKeyFileText := pksf.getPrivateKeyFileText()

		FileName := SaveFileDl.FileName()
		if strings.LastIndex(FileName, ".") < 0 ||
			strings.LastIndex(FileName, ".") <
				strings.LastIndex(FileName, "\\") {
			FileName += ".dat"
		}

		ioutil.WriteFile(FileName, []byte(PrivateKeyFileText), 0644)
		pksf.XgqA.Nickname = pksf.EditNickname.Text()
		VCI.AddAccount(&pksf.XgqA)
	} else {
		pksf.ClearPrivateKeyData()
	}

	Form.Close()
}

func (pksf *PrivateKeySaveForm) getPrivateKeyFileText() string {
	sc := pksf.sc
	StrAll := ""

	var bytes []byte
	var StrEn string

	Password := pksf.EditPassword.Text()
	pksf.EditPassword.SetText("")
	Nickname := pksf.EditNickname.Text()
	// pksf.EditNickname.SetText("")

	md := sha512.New()
	md.Write([]byte(Password))
	BytesKey := md.Sum(nil)[:32] // and compute ChopMD-256(SHA-512),

	StrAll += Nickname + "\r\n"
	StrAll += strconv.Itoa(pksf.curveType) + "\r\n"
	StrAll += pksf.addr + "\r\n"
	StrAll += pksf.x.String() + "\r\n"
	StrAll += pksf.y.String() + "\r\n"
	StrAll += pksf.d.String() + "\r\n"

	bytes, _ = sc.AesEncrypt([]byte(StrAll), BytesKey)
	StrEn = string(bytes[:])

	pksf.ClearPrivateKeyData()
	Password = ""
	Nickname = ""
	StrAll = ""
	BytesKey = []byte("0")

	return StrEn
}

func (pksf *PrivateKeySaveForm) ClearPrivateKeyData() {
	pksf.addr = ""
	pksf.x.SetBytes([]byte("0"))
	pksf.y.SetBytes([]byte("0"))
	pksf.d.SetBytes([]byte("0"))
}
